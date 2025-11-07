/**
 * S3: 移動指示 + ホラー演出
 *
 * - 2階診査室への導線
 * - カメラにホラー用フィルター（色味+ノイズ）を重ねる
 * - ランダムで影・視線・虚な人のシルエットを表示
 * - SE：低音の環境音 + 「こっちに来てはならない」
 */

'use client';

import { useState, useEffect } from 'react';
import styles from './page.module.css';
import Button from '@/components/ui/Button';
import CameraPreview from '@/components/camera/CameraPreview';
import HorrorOverlay from '@/components/horror/HorrorOverlay';
import GlitchText from '@/components/horror/GlitchText';
import AudioPlayer from '@/components/audio/AudioPlayer';
import { useGameFlow } from '@/hooks/useGameFlow';

export default function S3Page() {
  const { transitionTo } = useGameFlow();
  const [showShadows, setShowShadows] = useState(false);
  const [instruction, setInstruction] = useState('');

  // ホラー演出のタイミング制御
  useEffect(() => {
    // 3秒後に影を表示
    const shadowTimer = setTimeout(() => {
      setShowShadows(true);
    }, 3000);

    // 5秒後に指示を表示
    const instructionTimer = setTimeout(() => {
      setInstruction('2階の診査室へお越しください');
    }, 5000);

    // ランダムに影を点滅
    const flickerInterval = setInterval(() => {
      setShowShadows((prev) => !prev);
    }, 8000);

    return () => {
      clearTimeout(shadowTimer);
      clearTimeout(instructionTimer);
      clearInterval(flickerInterval);
    };
  }, []);

  const handleProceed = async () => {
    await transitionTo('s4');
  };

  return (
    <div className={styles.container}>
      {/* カメラプレビュー + ホラーオーバーレイ */}
      <div className={styles.cameraBackground}>
        <HorrorOverlay
          enableNoise
          enableVignette
          enableSepia
          enableRedTint
          noiseIntensity={0.3}
          showShadows={showShadows}
        >
          <CameraPreview facingMode="environment" />
        </HorrorOverlay>
      </div>

      {/* 環境音・ホラーSE */}
      <AudioPlayer
        src="/audio/horror_ambient.mp3"
        loop
        volume={0.3}
        autoPlay
      />

      {/* コンテンツ */}
      <div className={styles.content}>
        {/* 警告メッセージ（グリッチ効果） */}
        <div className={styles.warningBox}>
          <GlitchText intensity="high">
            こっちに来てはならない
          </GlitchText>
        </div>

        {/* 移動指示 */}
        {instruction && (
          <div className={styles.instructionBox}>
            <p className={styles.instruction}>{instruction}</p>

            <div className={styles.buttonContainer}>
              <Button
                onClick={handleProceed}
                variant="primary"
                size="large"
                fullWidth
              >
                診査室へ向かう
              </Button>
            </div>
          </div>
        )}

        {/* ヒント（下部固定） */}
        <div className={styles.hint}>
          <p>※ カメラを動かして周囲を確認してください</p>
        </div>
      </div>
    </div>
  );
}
