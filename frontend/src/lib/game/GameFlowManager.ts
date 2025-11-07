/**
 * GameFlowManager - ゲームフロー全体を管理する集中型マネージャー
 *
 * 責務:
 * - セッション初期化・管理
 * - シーン遷移制御（バリデーション付き）
 * - ブラウザバック完全ブロック
 * - エラーハンドリング
 * - LocalStorage ベースの状態管理
 */

import browserHistoryBlocker from './browser-history-blocker';
import {
  type SceneId,
  validateTransition,
  isValidSceneId,
  getSceneRoute,
  TransitionError,
} from './scene-transitions';
import {
  type GameSession,
  generateSessionId,
  calculateExpiryDate,
  getStoredSessionId,
  storeSessionId,
  clearStoredSessionId,
  fetchSession,
  registerSession,
  isSessionValid,
} from './session';

// LocalStorageキー
const CURRENT_SCENE_KEY = 'game_current_scene';

// ゲームエラー型定義
export interface GameError {
  code: string;
  message: string;
  userMessage: string;
  recoverable: boolean;
}

// イベントハンドラー型
export type ErrorHandler = (error: GameError) => void;
export type SceneChangeHandler = (sceneId: SceneId) => void;

/**
 * GameFlowManager クラス (Singleton)
 */
class GameFlowManager {
  private session: GameSession | null = null;
  private currentScene: SceneId = 's0';
  private isInitialized = false;
  private errorHandlers: ErrorHandler[] = [];
  private sceneChangeHandlers: SceneChangeHandler[] = [];

  /**
   * 初期化（アプリ起動時に1回だけ実行）
   */
  async initialize(): Promise<void> {
    if (this.isInitialized) {
      console.debug('[GameFlowManager] すでに初期化済みです');
      return;
    }

    try {
      // セッション復元または新規作成
      await this.restoreOrCreateSession();

      // 現在シーンを復元（LocalStorage）
      this.restoreCurrentScene();

      // ブラウザバックを有効化
      browserHistoryBlocker.enable();

      this.isInitialized = true;
      console.debug(
        '[GameFlowManager] 初期化完了',
        `Session: ${this.session?.sessionId}`,
        `Scene: ${this.currentScene}`
      );
    } catch (error) {
      console.error('[GameFlowManager] 初期化エラー:', error);
      this.handleError({
        code: 'INIT_ERROR',
        message: String(error),
        userMessage: 'ゲームの初期化に失敗しました',
        recoverable: false,
      });
    }
  }

  /**
   * セッションの復元または新規作成
   */
  private async restoreOrCreateSession(): Promise<void> {
    const storedSessionId = getStoredSessionId();

    if (storedSessionId) {
      // 既存セッションを復元
      const session = await fetchSession(storedSessionId);

      if (session && isSessionValid(session.expiresAt)) {
        this.session = session;
        console.debug('[GameFlowManager] セッション復元:', storedSessionId);
        return;
      }

      // 無効なセッションは削除
      clearStoredSessionId();
      console.debug('[GameFlowManager] 期限切れセッションを削除');
    }

    // 新規セッション作成
    const sessionId = generateSessionId();
    const createdAt = new Date().toISOString();
    const expiresAt = calculateExpiryDate().toISOString();

    this.session = {
      sessionId,
      createdAt,
      expiresAt,
      currentScene: 's0',
    };

    storeSessionId(sessionId);

    // サーバーに登録
    const registered = await registerSession(sessionId);
    if (!registered) {
      console.warn('[GameFlowManager] セッション登録に失敗（続行）');
    }

    console.debug('[GameFlowManager] 新規セッション作成:', sessionId);
  }

  /**
   * 現在シーンをLocalStorageから復元
   */
  private restoreCurrentScene(): void {
    if (typeof window === 'undefined') return;

    const stored = localStorage.getItem(CURRENT_SCENE_KEY);
    if (stored && isValidSceneId(stored)) {
      this.currentScene = stored;
      console.debug('[GameFlowManager] シーン復元:', stored);
    } else {
      // デフォルトはs0
      this.currentScene = 's0';
      this.saveCurrentScene();
    }
  }

  /**
   * 現在シーンをLocalStorageに保存
   */
  private saveCurrentScene(): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem(CURRENT_SCENE_KEY, this.currentScene);
  }

  /**
   * 現在のセッション取得
   */
  getSession(): GameSession | null {
    return this.session;
  }

  /**
   * 現在のシーンID取得
   */
  getCurrentScene(): SceneId {
    return this.currentScene;
  }

  /**
   * セッション有効性チェック
   */
  async validateSession(): Promise<boolean> {
    if (!this.session) {
      return false;
    }

    // ローカルの有効期限チェック
    if (!isSessionValid(this.session.expiresAt)) {
      this.handleSessionExpired();
      return false;
    }

    // サーバー側のセッション確認（オプション）
    const serverSession = await fetchSession(this.session.sessionId);
    if (!serverSession) {
      this.handleSessionExpired();
      return false;
    }

    return true;
  }

  /**
   * シーン遷移実行
   */
  async transitionTo(
    targetScene: SceneId,
    skipValidation: boolean = false
  ): Promise<boolean> {
    try {
      if (!this.isInitialized) {
        throw new Error('GameFlowManager が初期化されていません');
      }

      if (!this.session) {
        throw new Error('セッションが存在しません');
      }

      // バリデーションチェック
      if (!skipValidation) {
        const validation = await validateTransition(
          this.currentScene,
          targetScene,
          this.session.sessionId
        );

        if (!validation.valid && validation.error) {
          this.handleError({
            code: validation.error.code,
            message: validation.error.message,
            userMessage: validation.error.userMessage,
            recoverable: true,
          });
          return false;
        }
      }

      // シーン遷移実行
      const previousScene = this.currentScene;
      this.currentScene = targetScene;
      this.saveCurrentScene();

      // ハンドラー通知
      this.notifySceneChange(targetScene);

      console.debug(
        `[GameFlowManager] シーン遷移: ${previousScene} → ${targetScene}`
      );

      return true;
    } catch (error) {
      console.error('[GameFlowManager] 遷移エラー:', error);
      this.handleError({
        code: 'TRANSITION_ERROR',
        message: String(error),
        userMessage: 'シーン遷移に失敗しました',
        recoverable: true,
      });
      return false;
    }
  }

  /**
   * 指定シーンへ遷移可能かチェック
   */
  async canTransitionTo(targetScene: SceneId): Promise<boolean> {
    if (!this.session) {
      return false;
    }

    const validation = await validateTransition(
      this.currentScene,
      targetScene,
      this.session.sessionId
    );

    return validation.valid;
  }

  /**
   * セッション期限切れ処理
   */
  private handleSessionExpired(): void {
    this.handleError({
      code: 'SESSION_EXPIRED',
      message: 'Session has expired',
      userMessage: 'セッションの有効期限が切れました。最初からやり直してください。',
      recoverable: false,
    });

    // セッションをクリア
    clearStoredSessionId();
    this.session = null;

    // S0へリセット
    this.currentScene = 's0';
    this.saveCurrentScene();
  }

  /**
   * エラーハンドリング
   */
  private handleError(error: GameError): void {
    console.error('[GameFlowManager] Error:', error);

    // 登録されたエラーハンドラーに通知
    this.errorHandlers.forEach((handler) => {
      try {
        handler(error);
      } catch (e) {
        console.error('[GameFlowManager] エラーハンドラー実行エラー:', e);
      }
    });
  }

  /**
   * シーン変更通知
   */
  private notifySceneChange(sceneId: SceneId): void {
    this.sceneChangeHandlers.forEach((handler) => {
      try {
        handler(sceneId);
      } catch (e) {
        console.error('[GameFlowManager] シーン変更ハンドラー実行エラー:', e);
      }
    });
  }

  /**
   * エラーハンドラーを登録
   */
  onError(handler: ErrorHandler): () => void {
    this.errorHandlers.push(handler);

    // 登録解除関数を返す
    return () => {
      const index = this.errorHandlers.indexOf(handler);
      if (index !== -1) {
        this.errorHandlers.splice(index, 1);
      }
    };
  }

  /**
   * シーン変更ハンドラーを登録
   */
  onSceneChange(handler: SceneChangeHandler): () => void {
    this.sceneChangeHandlers.push(handler);

    // 登録解除関数を返す
    return () => {
      const index = this.sceneChangeHandlers.indexOf(handler);
      if (index !== -1) {
        this.sceneChangeHandlers.splice(index, 1);
      }
    };
  }

  /**
   * ゲーム終了（クリーンアップ）
   */
  dispose(): void {
    browserHistoryBlocker.disable();
    this.errorHandlers = [];
    this.sceneChangeHandlers = [];
    this.isInitialized = false;
    console.debug('[GameFlowManager] クリーンアップ完了');
  }

  /**
   * シーンをリセット（デバッグ用）
   */
  resetToScene(sceneId: SceneId): void {
    this.currentScene = sceneId;
    this.saveCurrentScene();
    console.debug(`[GameFlowManager] シーンリセット: ${sceneId}`);
  }
}

// シングルトンインスタンス
const gameFlowManager = new GameFlowManager();

export default gameFlowManager;
export { type GameError };
