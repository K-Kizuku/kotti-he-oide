/**
 * JumpScare - ジャンプスケア演出コンポーネント
 *
 * S6の誤答時などに使用
 * 画像とSEを組み合わせた驚き演出
 */

'use client';

import { useEffect, useState } from 'react';
import styles from './JumpScare.module.css';

interface JumpScareProps {
  /** ジャンプスケアを表示するか */
  isActive: boolean;
  /** 完了時のコールバック */
  onComplete?: () => void;
  /** 表示時間（ミリ秒） */
  duration?: number;
  /** 画像URL（オプション） */
  imageUrl?: string;
  /** テキスト表示（画像がない場合） */
  text?: string;
}

export default function JumpScare({
  isActive,
  onComplete,
  duration = 800,
  imageUrl,
  text = 'こっちに来てはならない',
}: JumpScareProps) {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    if (isActive) {
      setIsVisible(true);

      // 表示時間後に非表示にして完了コールバック
      const timer = setTimeout(() => {
        setIsVisible(false);
        onComplete?.();
      }, duration);

      return () => clearTimeout(timer);
    }
  }, [isActive, duration, onComplete]);

  if (!isVisible) return null;

  return (
    <div className={styles.container}>
      <div className={styles.flash} />

      {imageUrl ? (
        <img
          src={imageUrl}
          alt="ジャンプスケア"
          className={styles.image}
        />
      ) : (
        <div className={styles.textContainer}>
          <p className={styles.text}>{text}</p>
        </div>
      )}

      <div className={styles.glitchOverlay} />
    </div>
  );
}
