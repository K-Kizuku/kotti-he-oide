/**
 * useTimer - タイマー管理フック
 *
 * S5, S6, S8で使用
 * カウントダウンタイマーの状態管理
 */

import { useState, useEffect, useCallback, useRef } from 'react';

interface UseTimerOptions {
  /** 初期秒数 */
  initialSeconds: number;
  /** 自動開始するか */
  autoStart?: boolean;
  /** 完了時のコールバック */
  onComplete?: () => void;
  /** 警告閾値（秒） */
  warningThreshold?: number;
}

export function useTimer({
  initialSeconds,
  autoStart = false,
  onComplete,
  warningThreshold = 60,
}: UseTimerOptions) {
  const [seconds, setSeconds] = useState(initialSeconds);
  const [isRunning, setIsRunning] = useState(autoStart);
  const [isWarning, setIsWarning] = useState(false);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  // タイマー開始
  const start = useCallback(() => {
    setIsRunning(true);
  }, []);

  // タイマー停止
  const pause = useCallback(() => {
    setIsRunning(false);
  }, []);

  // タイマーリセット
  const reset = useCallback((newSeconds?: number) => {
    setSeconds(newSeconds ?? initialSeconds);
    setIsRunning(false);
    setIsWarning(false);
  }, [initialSeconds]);

  // 秒数追加
  const addSeconds = useCallback((additionalSeconds: number) => {
    setSeconds((prev) => Math.max(0, prev + additionalSeconds));
  }, []);

  // タイマーロジック
  useEffect(() => {
    if (!isRunning) return;

    intervalRef.current = setInterval(() => {
      setSeconds((prev) => {
        const newValue = prev - 1;

        // 警告閾値チェック
        if (newValue <= warningThreshold && !isWarning) {
          setIsWarning(true);
        }

        // 完了チェック
        if (newValue <= 0) {
          setIsRunning(false);
          onComplete?.();
          return 0;
        }

        return newValue;
      });
    }, 1000);

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [isRunning, warningThreshold, isWarning, onComplete]);

  const minutes = Math.floor(seconds / 60);
  const secs = seconds % 60;
  const timeString = `${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
  const percentage = (seconds / initialSeconds) * 100;

  return {
    seconds,
    minutes,
    secs,
    timeString,
    percentage,
    isRunning,
    isWarning,
    start,
    pause,
    reset,
    addSeconds,
  };
}
