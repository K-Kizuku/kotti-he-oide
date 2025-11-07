'use client';

/**
 * useGameFlow - ゲームフロー管理用カスタムフック
 *
 * GameFlowContextにアクセスし、各コンポーネントでゲームフロー制御を行います
 */

import { useContext } from 'react';
import { GameFlowContext, type GameFlowContextValue } from '@/contexts/GameFlowContext';

/**
 * useGameFlow フック
 *
 * @throws {Error} GameFlowProviderの外で使用された場合
 * @returns {GameFlowContextValue} ゲームフロー管理のContext値
 *
 * @example
 * ```tsx
 * const { session, currentScene, transitionTo } = useGameFlow();
 *
 * const handleNext = async () => {
 *   const success = await transitionTo('s1');
 *   if (!success) {
 *     console.error('遷移に失敗しました');
 *   }
 * };
 * ```
 */
export function useGameFlow(): GameFlowContextValue {
  const context = useContext(GameFlowContext);

  if (!context) {
    throw new Error(
      'useGameFlow must be used within a GameFlowProvider. ' +
      'Ensure that your component is wrapped with <GameFlowProvider>.'
    );
  }

  return context;
}
