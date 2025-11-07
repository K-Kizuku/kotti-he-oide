/**
 * HorrorOverlay - ホラー演出オーバーレイコンポーネント
 *
 * ノイズ・グリッチ効果、ヴィネット効果、色調補正などを適用
 * S3以降のシーンで使用
 */

'use client';

import { ReactNode } from 'react';
import styles from './HorrorOverlay.module.css';

interface HorrorOverlayProps {
  /** 子要素（カメラプレビューやコンテンツ） */
  children: ReactNode;
  /** ノイズ効果を適用するか */
  enableNoise?: boolean;
  /** ヴィネット効果を適用するか */
  enableVignette?: boolean;
  /** セピア調フィルターを適用するか */
  enableSepia?: boolean;
  /** グリッチ効果を適用するか */
  enableGlitch?: boolean;
  /** 赤み強調（血の色） */
  enableRedTint?: boolean;
  /** ノイズの強度（0-1） */
  noiseIntensity?: number;
  /** 影・シルエットを表示するか */
  showShadows?: boolean;
}

export default function HorrorOverlay({
  children,
  enableNoise = true,
  enableVignette = true,
  enableSepia = false,
  enableGlitch = false,
  enableRedTint = false,
  noiseIntensity = 0.5,
  showShadows = false,
}: HorrorOverlayProps) {
  const overlayClassNames = [
    styles.container,
    enableNoise && styles.noise,
    enableVignette && styles.vignette,
    enableSepia && styles.sepia,
    enableGlitch && styles.glitch,
    enableRedTint && styles.redTint,
  ].filter(Boolean).join(' ');

  return (
    <div className={overlayClassNames}>
      {children}

      {/* ノイズレイヤー */}
      {enableNoise && (
        <div
          className={styles.noiseLayer}
          style={{ opacity: noiseIntensity }}
        />
      )}

      {/* ヴィネットレイヤー */}
      {enableVignette && (
        <div className={styles.vignetteLayer} />
      )}

      {/* 影・シルエットレイヤー */}
      {showShadows && (
        <div className={styles.shadowsLayer}>
          <div className={styles.shadow} style={{
            top: `${Math.random() * 60 + 20}%`,
            left: `${Math.random() * 60 + 20}%`,
          }} />
        </div>
      )}
    </div>
  );
}
