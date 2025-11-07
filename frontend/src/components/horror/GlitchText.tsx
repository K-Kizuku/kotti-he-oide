/**
 * GlitchText - グリッチテキストアニメーション
 *
 * ホラー演出用のテキストエフェクト
 * タイトルや警告文で使用
 */

'use client';

import { ReactNode } from 'react';
import styles from './GlitchText.module.css';

interface GlitchTextProps {
  /** 表示するテキスト */
  children: ReactNode;
  /** グリッチ効果の強度 */
  intensity?: 'low' | 'medium' | 'high';
  /** 常時アニメーションするか */
  continuous?: boolean;
  /** カスタムクラス名 */
  className?: string;
}

export default function GlitchText({
  children,
  intensity = 'medium',
  continuous = true,
  className = '',
}: GlitchTextProps) {
  const glitchClassNames = [
    styles.glitch,
    styles[intensity],
    continuous && styles.continuous,
    className,
  ].filter(Boolean).join(' ');

  return (
    <div className={glitchClassNames} data-text={children}>
      {children}
    </div>
  );
}
