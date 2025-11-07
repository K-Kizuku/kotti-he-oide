'use client';

/**
 * GameFlowContext - ゲームフロー管理のReact Context
 *
 * GameFlowManagerをReactコンポーネントから使用できるようにします
 */

import React, { createContext, useEffect, useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import gameFlowManager, { type GameError } from '@/lib/game/GameFlowManager';
import type { GameSession } from '@/lib/game/session';
import type { SceneId } from '@/lib/game/scene-transitions';
import { getSceneRoute } from '@/lib/game/scene-transitions';

// Context型定義
export interface GameFlowContextValue {
  // 状態
  session: GameSession | null;
  currentScene: SceneId;
  isInitialized: boolean;

  // メソッド
  transitionTo: (targetScene: SceneId, skipValidation?: boolean) => Promise<boolean>;
  canTransitionTo: (targetScene: SceneId) => Promise<boolean>;
  validateSession: () => Promise<boolean>;

  // エラー処理
  currentError: GameError | null;
  clearError: () => void;
}

// Contextの作成
export const GameFlowContext = createContext<GameFlowContextValue | null>(null);

// Provider Props
interface GameFlowProviderProps {
  children: React.ReactNode;
}

/**
 * GameFlowProvider - アプリ全体をラップするProvider
 */
export function GameFlowProvider({ children }: GameFlowProviderProps) {
  const router = useRouter();

  // 状態管理
  const [session, setSession] = useState<GameSession | null>(null);
  const [currentScene, setCurrentScene] = useState<SceneId>('s0');
  const [isInitialized, setIsInitialized] = useState(false);
  const [currentError, setCurrentError] = useState<GameError | null>(null);

  // 初期化
  useEffect(() => {
    const init = async () => {
      try {
        await gameFlowManager.initialize();

        const sessionData = gameFlowManager.getSession();
        const scene = gameFlowManager.getCurrentScene();

        setSession(sessionData);
        setCurrentScene(scene);
        setIsInitialized(true);

        console.debug('[GameFlowProvider] 初期化完了');
      } catch (error) {
        console.error('[GameFlowProvider] 初期化エラー:', error);
        setCurrentError({
          code: 'INIT_ERROR',
          message: String(error),
          userMessage: 'ゲームの初期化に失敗しました',
          recoverable: false,
        });
      }
    };

    init();

    // クリーンアップ
    return () => {
      gameFlowManager.dispose();
    };
  }, []);

  // エラーハンドラー登録
  useEffect(() => {
    const unsubscribe = gameFlowManager.onError((error) => {
      setCurrentError(error);

      // セッション期限切れの場合、S0へリダイレクト
      if (error.code === 'SESSION_EXPIRED') {
        router.push('/game/s0');
      }
    });

    return unsubscribe;
  }, [router]);

  // シーン変更ハンドラー登録
  useEffect(() => {
    const unsubscribe = gameFlowManager.onSceneChange((sceneId) => {
      setCurrentScene(sceneId);

      // Next.jsルーターで実際のページ遷移を実行
      const route = getSceneRoute(sceneId);
      router.push(route);
    });

    return unsubscribe;
  }, [router]);

  // シーン遷移メソッド
  const transitionTo = useCallback(
    async (targetScene: SceneId, skipValidation: boolean = false): Promise<boolean> => {
      const success = await gameFlowManager.transitionTo(targetScene, skipValidation);
      if (success) {
        setCurrentScene(targetScene);
      }
      return success;
    },
    []
  );

  // 遷移可能性チェック
  const canTransitionTo = useCallback(
    async (targetScene: SceneId): Promise<boolean> => {
      return await gameFlowManager.canTransitionTo(targetScene);
    },
    []
  );

  // セッション検証
  const validateSession = useCallback(async (): Promise<boolean> => {
    return await gameFlowManager.validateSession();
  }, []);

  // エラークリア
  const clearError = useCallback(() => {
    setCurrentError(null);
  }, []);

  // Context値
  const contextValue: GameFlowContextValue = {
    session,
    currentScene,
    isInitialized,
    transitionTo,
    canTransitionTo,
    validateSession,
    currentError,
    clearError,
  };

  return (
    <GameFlowContext.Provider value={contextValue}>
      {children}
    </GameFlowContext.Provider>
  );
}
