/**
 * useSession - セッション管理フック
 *
 * ゲームセッションの初期化・管理
 */

import { useState, useEffect } from 'react';
import {
  type GameSession,
  initializeSession,
  registerSession,
  clearStoredSessionId,
} from '@/lib/game/session';

interface UseSessionReturn {
  session: GameSession | null;
  isLoading: boolean;
  error: Error | null;
  refreshSession: () => Promise<void>;
  endSession: () => void;
}

export function useSession(): UseSessionReturn {
  const [session, setSession] = useState<GameSession | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const loadSession = async () => {
    try {
      setIsLoading(true);
      setError(null);

      const gameSession = await initializeSession();
      setSession(gameSession);

      // サーバーに登録（既存の場合は無視される）
      await registerSession(gameSession.sessionId);
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to initialize session');
      setError(error);
      console.error('Session initialization error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadSession();
  }, []);

  const refreshSession = async () => {
    await loadSession();
  };

  const endSession = () => {
    clearStoredSessionId();
    setSession(null);
  };

  return {
    session,
    isLoading,
    error,
    refreshSession,
    endSession,
  };
}
