/**
 * session.ts - セッション管理ユーティリティ
 *
 * ゲームセッションの作成、取得、保存を管理
 * LocalStorageを使用したクライアント側の永続化
 */

const SESSION_KEY = 'game_session_id';
const SESSION_EXPIRY_MINUTES = 60;

export interface GameSession {
  sessionId: string;
  createdAt: string;
  expiresAt: string;
  currentScene: string;
}

/**
 * 新しいセッションIDを生成（UUID v4形式）
 */
export function generateSessionId(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

/**
 * セッションの有効期限を計算
 */
export function calculateExpiryDate(minutes: number = SESSION_EXPIRY_MINUTES): Date {
  const now = new Date();
  return new Date(now.getTime() + minutes * 60 * 1000);
}

/**
 * セッションが有効かチェック
 */
export function isSessionValid(expiresAt: string): boolean {
  return new Date(expiresAt) > new Date();
}

/**
 * LocalStorageからセッションIDを取得
 */
export function getStoredSessionId(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(SESSION_KEY);
}

/**
 * LocalStorageにセッションIDを保存
 */
export function storeSessionId(sessionId: string): void {
  if (typeof window === 'undefined') return;
  localStorage.setItem(SESSION_KEY, sessionId);
}

/**
 * LocalStorageからセッションIDを削除
 */
export function clearStoredSessionId(): void {
  if (typeof window === 'undefined') return;
  localStorage.removeItem(SESSION_KEY);
}

/**
 * 新しいセッションを作成
 */
export function createNewSession(currentScene: string = 's0'): GameSession {
  const sessionId = generateSessionId();
  const createdAt = new Date().toISOString();
  const expiresAt = calculateExpiryDate().toISOString();

  const session: GameSession = {
    sessionId,
    createdAt,
    expiresAt,
    currentScene,
  };

  storeSessionId(sessionId);
  return session;
}

/**
 * セッション情報をAPIから取得
 */
export async function fetchSession(sessionId: string): Promise<GameSession | null> {
  try {
    const response = await fetch(`/api/session/${sessionId}`);

    if (!response.ok) {
      if (response.status === 404) {
        return null;
      }
      throw new Error(`Failed to fetch session: ${response.statusText}`);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching session:', error);
    return null;
  }
}

/**
 * セッションを初期化（既存セッション復元 or 新規作成）
 */
export async function initializeSession(): Promise<GameSession> {
  const storedSessionId = getStoredSessionId();

  if (storedSessionId) {
    const session = await fetchSession(storedSessionId);

    if (session && isSessionValid(session.expiresAt)) {
      return session;
    }

    // 無効なセッションは削除
    clearStoredSessionId();
  }

  // 新規セッション作成
  return createNewSession();
}

/**
 * セッションIDをサーバーに登録
 */
export async function registerSession(sessionId: string): Promise<boolean> {
  try {
    const response = await fetch('/api/session', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        session_id: sessionId,
      }),
    });

    return response.ok;
  } catch (error) {
    console.error('Error registering session:', error);
    return false;
  }
}
