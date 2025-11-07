/**
 * Timer - カウントダウンタイマーコンポーネント
 *
 * S5, S6, S8で使用
 * 残り時間を視覚的に表示
 */

'use client';

import { useEffect, useState } from 'react';
import styles from './Timer.module.css';

interface TimerProps {
  /** 開始時間（秒） */
  initialSeconds: number;
  /** タイマー終了時のコールバック */
  onComplete?: () => void;
  /** 警告閾値（秒） */
  warningThreshold?: number;
  /** 警告中かどうか */
  isWarning?: boolean;
  /** 大きいサイズで表示するか */
  large?: boolean;
  /** 常に表示（上部固定）するか */
  fixed?: boolean;
}

export default function Timer({
  initialSeconds,
  onComplete,
  warningThreshold = 60,
  large = false,
  fixed = false,
}: TimerProps) {
  const [seconds, setSeconds] = useState(initialSeconds);
  const [isWarning, setIsWarning] = useState(false);

  useEffect(() => {
    if (seconds <= 0) {
      onComplete?.();
      return;
    }

    const timer = setInterval(() => {
      setSeconds((prev) => {
        const newValue = prev - 1;
        if (newValue <= warningThreshold) {
          setIsWarning(true);
        }
        if (newValue <= 0) {
          clearInterval(timer);
          onComplete?.();
        }
        return newValue;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [seconds, onComplete, warningThreshold]);

  const minutes = Math.floor(seconds / 60);
  const secs = seconds % 60;
  const timeString = `${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;

  const percentage = (seconds / initialSeconds) * 100;

  const classNames = [
    styles.timer,
    large && styles.large,
    fixed && styles.fixed,
    isWarning && styles.warning,
  ].filter(Boolean).join(' ');

  return (
    <div className={classNames}>
      <div className={styles.display}>
        <span className={styles.time}>{timeString}</span>
        {isWarning && <span className={styles.warningIcon}>⚠️</span>}
      </div>
      <div className={styles.progressBar}>
        <div
          className={styles.progressFill}
          style={{ width: `${percentage}%` }}
        />
      </div>
    </div>
  );
}
