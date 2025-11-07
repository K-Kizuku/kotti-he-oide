/**
 * シーン遷移ルールと検証ロジック
 * 各シーンへの遷移条件を定義し、バリデーションを実行します
 */

import { getS6Progress, getAnswers, type S6Progress } from './api';

// シーンタイプ定義
export type SceneId =
  | 's0'
  | 's1'
  | 's2'
  | 's3'
  | 's4'
  | 's5'
  | 's6'
  | 's7'
  | 's8'
  | 's9'
  | 'gameover';

// 遷移検証エラー
export class TransitionError extends Error {
  constructor(
    message: string,
    public code: string,
    public userMessage: string
  ) {
    super(message);
    this.name = 'TransitionError';
  }
}

// シーン遷移グラフ定義（どのシーンからどのシーンへ遷移可能か）
const SCENE_GRAPH: Record<SceneId, SceneId[]> = {
  s0: ['s1'], // S0 → S1のみ
  s1: ['s2'], // S1 → S2のみ
  s2: ['s3'], // S2 → S3のみ
  s3: ['s4'], // S3 → S4のみ
  s4: ['s5'], // S4 → S5のみ
  s5: ['s6'], // S5 → S6のみ
  s6: ['s7', 'gameover'], // S6 → S7 または gameover（タイムアウト）
  s7: ['s8'], // S7 → S8のみ
  s8: ['s9', 'gameover'], // S8 → S9 または gameover（タイムアウト）
  s9: ['s0'], // S9 → S0（ゲーム終了後トップへ）
  gameover: ['s0'], // gameover → S0（リスタート）
};

/**
 * 基本的な遷移可能性チェック（グラフベース）
 */
export function isTransitionAllowed(from: SceneId, to: SceneId): boolean {
  const allowedTargets = SCENE_GRAPH[from];
  return allowedTargets.includes(to);
}

/**
 * S6 → S7 の遷移条件検証
 * 条件: 5箇所すべてでピースを取得済み（verified && correct）
 */
async function validateS6ToS7Transition(
  sessionId: string
): Promise<{ valid: boolean; error?: TransitionError }> {
  try {
    const progress = await getS6Progress(sessionId);

    if (progress.length !== 5) {
      return {
        valid: false,
        error: new TransitionError(
          'S6 progress incomplete: not all places verified',
          'S6_INCOMPLETE',
          '5箇所すべてのピースを集めてください'
        ),
      };
    }

    // 全ての場所で verified=true かつ correct=true を確認
    const allCompleted = progress.every((p) => p.verified && p.correct);

    if (!allCompleted) {
      const completedCount = progress.filter((p) => p.verified && p.correct).length;
      return {
        valid: false,
        error: new TransitionError(
          `S6 progress incomplete: ${completedCount}/5 pieces obtained`,
          'S6_INCOMPLETE',
          `まだ全てのピースが揃っていません（${completedCount}/5）`
        ),
      };
    }

    return { valid: true };
  } catch (error) {
    console.error('Failed to validate S6 → S7 transition:', error);
    return {
      valid: false,
      error: new TransitionError(
        'Failed to check S6 progress',
        'VALIDATION_ERROR',
        '進捗状況の確認に失敗しました'
      ),
    };
  }
}

/**
 * S7 → S8 の遷移条件検証
 * 条件: S4で答えた「人生の最期に達成したいこと」を正確に再入力済み
 * 注: この検証は通常クライアント側で行われるため、ここでは基本チェックのみ
 */
async function validateS7ToS8Transition(
  sessionId: string
): Promise<{ valid: boolean; error?: TransitionError }> {
  try {
    const answers = await getAnswers(sessionId);

    // question_id="9" が「人生の最期に達成したいこと」と仮定
    // （実際のquestion_idは constants.ts の INTROSPECTION_QUESTIONS を参照）
    const lifeGoalAnswer = answers.find(
      (a) => a.question_id === '9' || a.question_id === 'q9_life_goal'
    );

    if (!lifeGoalAnswer) {
      return {
        valid: false,
        error: new TransitionError(
          'Life goal answer not found',
          'ANSWER_NOT_FOUND',
          'S4で「人生の最期に達成したいこと」に回答していません'
        ),
      };
    }

    // 実際の検証はS7ページ内で行われるため、ここではtrueを返す
    return { valid: true };
  } catch (error) {
    console.error('Failed to validate S7 → S8 transition:', error);
    return {
      valid: false,
      error: new TransitionError(
        'Failed to check answers',
        'VALIDATION_ERROR',
        '回答の確認に失敗しました'
      ),
    };
  }
}

/**
 * 汎用シーン遷移検証
 */
export async function validateTransition(
  from: SceneId,
  to: SceneId,
  sessionId: string
): Promise<{ valid: boolean; error?: TransitionError }> {
  // 1. グラフベースの基本チェック
  if (!isTransitionAllowed(from, to)) {
    return {
      valid: false,
      error: new TransitionError(
        `Transition ${from} → ${to} not allowed`,
        'INVALID_TRANSITION',
        'この遷移は許可されていません'
      ),
    };
  }

  // 2. 特殊条件のチェック
  if (from === 's6' && to === 's7') {
    return await validateS6ToS7Transition(sessionId);
  }

  if (from === 's7' && to === 's8') {
    return await validateS7ToS8Transition(sessionId);
  }

  // その他の遷移は条件なし
  return { valid: true };
}

/**
 * シーンIDの妥当性チェック
 */
export function isValidSceneId(sceneId: string): sceneId is SceneId {
  return [
    's0',
    's1',
    's2',
    's3',
    's4',
    's5',
    's6',
    's7',
    's8',
    's9',
    'gameover',
  ].includes(sceneId);
}

/**
 * シーンIDからルートパスを生成
 */
export function getSceneRoute(sceneId: SceneId): string {
  if (sceneId === 'gameover') {
    return '/game/gameover';
  }
  return `/game/${sceneId}`;
}

/**
 * 現在のパスからシーンIDを抽出
 */
export function extractSceneIdFromPath(path: string): SceneId | null {
  const match = path.match(/\/game\/(s\d+|gameover)/);
  if (!match) {
    return null;
  }

  const sceneId = match[1];
  return isValidSceneId(sceneId) ? sceneId : null;
}
