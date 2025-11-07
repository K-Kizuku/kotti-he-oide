/**
 * useGameState - ゲーム状態管理フック
 *
 * シーン遷移、進行状況の管理
 */

import { useState, useCallback } from 'react';

export type SceneId =
  | 's0' // 起動・注意書き
  | 's1' // 1942年パート
  | 's2' // お気に入りの場所
  | 's3' // 移動指示
  | 's4' // 診査室質問
  | 's5' // 死亡届受理
  | 's6' // 存在証明書探索
  | 's7' // 2002年パート
  | 's8' // メインホール探索
  | 's9' // メッセージ刻み
  | 'gameover'; // ゲームオーバー

interface UseGameStateOptions {
  sessionId?: string;
  initialScene?: SceneId;
}

interface UseGameStateReturn {
  currentScene: SceneId;
  goToScene: (scene: SceneId) => void;
  nextScene: () => void;
  previousScene: () => void;
  isFirstScene: boolean;
  isLastScene: boolean;
}

const SCENE_ORDER: SceneId[] = [
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
];

export function useGameState({
  sessionId,
  initialScene = 's0',
}: UseGameStateOptions = {}): UseGameStateReturn {
  const [currentScene, setCurrentScene] = useState<SceneId>(initialScene);

  const goToScene = useCallback(
    (scene: SceneId) => {
      setCurrentScene(scene);

      // サーバーにシーン更新を通知
      if (sessionId) {
        fetch(`/api/session/${sessionId}/scene`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            current_scene: scene,
          }),
        }).catch((err) => {
          console.error('Failed to update scene:', err);
        });
      }
    },
    [sessionId]
  );

  const nextScene = useCallback(() => {
    const currentIndex = SCENE_ORDER.indexOf(currentScene);
    if (currentIndex < SCENE_ORDER.length - 1) {
      goToScene(SCENE_ORDER[currentIndex + 1]);
    }
  }, [currentScene, goToScene]);

  const previousScene = useCallback(() => {
    const currentIndex = SCENE_ORDER.indexOf(currentScene);
    if (currentIndex > 0) {
      goToScene(SCENE_ORDER[currentIndex - 1]);
    }
  }, [currentScene, goToScene]);

  const isFirstScene = currentScene === SCENE_ORDER[0];
  const isLastScene = currentScene === SCENE_ORDER[SCENE_ORDER.length - 1];

  return {
    currentScene,
    goToScene,
    nextScene,
    previousScene,
    isFirstScene,
    isLastScene,
  };
}
