"use client";

import styles from "./BackgroundEffects.module.css";

type BackgroundEffectsProps = {
  enableNoise?: boolean;
  enableScanlines?: boolean;
  enableVignette?: boolean;
};

export default function BackgroundEffects({
  enableNoise = true,
  enableScanlines = true,
  enableVignette = true,
}: BackgroundEffectsProps) {
  return (
    <div className={styles.effectsContainer} aria-hidden>
      {/* SVG ノイズフィルター */}
      <svg className={styles.hiddenSvg} width="0" height="0">
        <filter id="noiseFilter">
          <feTurbulence type="fractalNoise" baseFrequency="0.8" numOctaves={2} stitchTiles="stitch" />
          <feColorMatrix type="saturate" values="0" />
        </filter>
      </svg>

      {enableNoise && <div className={styles.noise} />}
      {enableScanlines && <div className={styles.scanlines} />}
      {enableVignette && <div className={styles.vignette} />}
    </div>
  );
}

