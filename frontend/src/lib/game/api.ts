/**
 * api.ts - ゲームAPI呼び出しユーティリティ
 *
 * サーバーとの通信を管理
 */

/**
 * S4の回答を保存
 */
export interface SaveAnswerRequest {
  question_id: string;
  answer_text: string;
}

export async function saveAnswer(
  sessionId: string,
  answer: SaveAnswerRequest
): Promise<boolean> {
  try {
    const response = await fetch(`/api/session/${sessionId}/answers`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(answer),
    });

    return response.ok;
  } catch (error) {
    console.error('Failed to save answer:', error);
    return false;
  }
}

/**
 * S4の保存済み回答を取得
 */
export interface SavedAnswer {
  question_id: string;
  answer_text: string;
  answered_at: string;
}

export async function getAnswers(sessionId: string): Promise<SavedAnswer[]> {
  try {
    const response = await fetch(`/api/session/${sessionId}/answers`);

    if (!response.ok) {
      throw new Error(`Failed to fetch answers: ${response.statusText}`);
    }

    const data = await response.json();
    return data.answers || [];
  } catch (error) {
    console.error('Failed to fetch answers:', error);
    return [];
  }
}

/**
 * S6探索開始
 */
export async function startS6Exploration(sessionId: string): Promise<boolean> {
  try {
    const response = await fetch(`/api/session/${sessionId}/s6/start`, {
      method: 'POST',
    });

    return response.ok;
  } catch (error) {
    console.error('Failed to start S6:', error);
    return false;
  }
}

/**
 * S6場所到達判定（画像類似度）
 */
export interface VerifyLocationRequest {
  place_id: string;
  image: Blob;
}

export interface VerifyLocationResponse {
  verified: boolean;
  similarity: number;
}

export async function verifyLocation(
  sessionId: string,
  request: VerifyLocationRequest
): Promise<VerifyLocationResponse> {
  try {
    const formData = new FormData();
    formData.append('place_id', request.place_id);
    formData.append('image', request.image, 'photo.jpg');

    const response = await fetch(
      `/api/session/${sessionId}/s6/verify-location`,
      {
        method: 'POST',
        body: formData,
      }
    );

    if (!response.ok) {
      throw new Error(`Failed to verify location: ${response.statusText}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to verify location:', error);
    return { verified: false, similarity: 0 };
  }
}

/**
 * S6クイズ取得
 */
export interface QuizOption {
  id: string;
  label: string;
}

export interface Quiz {
  quiz_id: string;
  place_id: string;
  question: string;
  options: QuizOption[];
  correct_answer_id: string;
}

export async function getQuiz(
  sessionId: string,
  placeId: string
): Promise<Quiz | null> {
  try {
    const response = await fetch(
      `/api/session/${sessionId}/s6/quiz/${placeId}`
    );

    if (!response.ok) {
      throw new Error(`Failed to fetch quiz: ${response.statusText}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to fetch quiz:', error);
    return null;
  }
}

/**
 * S6クイズ回答送信
 */
export interface AnswerQuizRequest {
  quiz_id: string;
  selected_answer_id: string;
}

export interface AnswerQuizResponse {
  correct: boolean;
  piece_obtained: boolean;
}

export async function answerQuiz(
  sessionId: string,
  request: AnswerQuizRequest
): Promise<AnswerQuizResponse> {
  try {
    const response = await fetch(`/api/session/${sessionId}/s6/answer`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      throw new Error(`Failed to answer quiz: ${response.statusText}`);
    }

    return await response.json();
  } catch (error) {
    console.error('Failed to answer quiz:', error);
    return { correct: false, piece_obtained: false };
  }
}

/**
 * S6進捗状況取得
 */
export interface S6Progress {
  place_id: string;
  verified: boolean;
  answered: boolean;
  correct: boolean;
}

export async function getS6Progress(
  sessionId: string
): Promise<S6Progress[]> {
  try {
    const response = await fetch(`/api/session/${sessionId}/s6/progress`);

    if (!response.ok) {
      throw new Error(`Failed to fetch S6 progress: ${response.statusText}`);
    }

    const data = await response.json();
    return data.progress || [];
  } catch (error) {
    console.error('Failed to fetch S6 progress:', error);
    return [];
  }
}

/**
 * VOICEVOX音声生成
 */
export interface GenerateVoiceRequest {
  text: string;
  speaker_id?: number;
}

export interface GenerateVoiceResponse {
  audio_url: string;
}

export async function generateVoice(
  request: GenerateVoiceRequest
): Promise<string | null> {
  try {
    const response = await fetch('/api/voice/generate', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        text: request.text,
        speaker_id: request.speaker_id || 11, // 青山龍星(しっとり)
      }),
    });

    if (!response.ok) {
      throw new Error(`Failed to generate voice: ${response.statusText}`);
    }

    const data: GenerateVoiceResponse = await response.json();
    return data.audio_url;
  } catch (error) {
    console.error('Failed to generate voice:', error);
    return null;
  }
}

/**
 * S9最終メッセージ保存
 */
export interface SaveMessageRequest {
  message_text: string;
  place_id: string;
}

export async function saveMessage(
  sessionId: string,
  request: SaveMessageRequest
): Promise<boolean> {
  try {
    const response = await fetch(`/api/session/${sessionId}/message`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    return response.ok;
  } catch (error) {
    console.error('Failed to save message:', error);
    return false;
  }
}

/**
 * 過去プレイヤーのメッセージ取得
 */
export interface PlayerMessage {
  message_id: string;
  message_text: string;
  place_id: string;
  created_at: string;
}

export async function getMessages(
  placeId?: string
): Promise<PlayerMessage[]> {
  try {
    const url = placeId
      ? `/api/messages?place_id=${placeId}`
      : '/api/messages';

    const response = await fetch(url);

    if (!response.ok) {
      throw new Error(`Failed to fetch messages: ${response.statusText}`);
    }

    const data = await response.json();
    return data.messages || [];
  } catch (error) {
    console.error('Failed to fetch messages:', error);
    return [];
  }
}
