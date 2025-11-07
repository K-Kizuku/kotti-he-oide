/**
 * CameraCapture - 撮影機能付きカメラコンポーネント
 *
 * S6の画像類似度判定で使用
 * 撮影ボタンを含む
 */

'use client';

import { useRef, useState } from 'react';
import styles from './CameraCapture.module.css';
import CameraPreview from './CameraPreview';
import Button from '../ui/Button';
import { FilterId } from '@/_archive/app-pages/camera-filters/filters';

interface CameraCaptureProps {
  /** フィルターID */
  filterId?: FilterId | null;
  /** カメラ方向 */
  facingMode?: 'user' | 'environment';
  /** 撮影完了ハンドラ */
  onCapture: (blob: Blob) => void;
  /** キャンセルハンドラ */
  onCancel?: () => void;
  /** 撮影ボタンのラベル */
  captureButtonLabel?: string;
  /** キャンセルボタンのラベル */
  cancelButtonLabel?: string;
}

export default function CameraCapture({
  filterId = 'horror',
  facingMode = 'environment',
  onCapture,
  onCancel,
  captureButtonLabel = '撮影する',
  cancelButtonLabel = 'この場所にいることにする',
}: CameraCaptureProps) {
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isReady, setIsReady] = useState(false);
  const [isCapturing, setIsCapturing] = useState(false);

  const handleCapture = async () => {
    if (!videoRef.current || !canvasRef.current) return;

    setIsCapturing(true);

    try {
      const video = videoRef.current;
      const canvas = canvasRef.current;
      const ctx = canvas.getContext('2d');

      if (!ctx) {
        throw new Error('Canvas context not available');
      }

      // 現在のビデオフレームをキャプチャ
      canvas.width = video.videoWidth;
      canvas.height = video.videoHeight;
      ctx.drawImage(video, 0, 0, canvas.width, canvas.height);

      // Blobに変換
      canvas.toBlob(
        (blob) => {
          if (blob) {
            onCapture(blob);
          }
          setIsCapturing(false);
        },
        'image/jpeg',
        0.8
      );
    } catch (error) {
      console.error('撮影エラー:', error);
      setIsCapturing(false);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.preview}>
        <CameraPreview
          filterId={filterId}
          facingMode={facingMode}
          onReady={(stream) => {
            setIsReady(true);
            // videoRefに参照を保存
            const video = document.querySelector('video');
            if (video) {
              videoRef.current = video;
            }
          }}
          onError={(error) => {
            console.error('カメラエラー:', error);
          }}
        />
      </div>

      {/* 非表示のキャンバス（撮影用） */}
      <canvas ref={canvasRef} style={{ display: 'none' }} />

      <div className={styles.controls}>
        <Button
          onClick={handleCapture}
          disabled={!isReady || isCapturing}
          isLoading={isCapturing}
          variant="primary"
          size="large"
          fullWidth
        >
          {captureButtonLabel}
        </Button>

        {onCancel && (
          <Button
            onClick={onCancel}
            variant="ghost"
            size="medium"
            fullWidth
          >
            {cancelButtonLabel}
          </Button>
        )}
      </div>
    </div>
  );
}
