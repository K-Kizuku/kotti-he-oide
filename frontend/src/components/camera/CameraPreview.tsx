/**
 * CameraPreview - カメラプレビューコンポーネント
 *
 * カメラ映像をリアルタイムで表示
 * フィルター適用可能
 */

'use client';

import { useEffect, useRef, useState } from 'react';
import styles from './CameraPreview.module.css';
import { applyFilter, FilterId } from '@/_archive/app-pages/camera-filters/filters';

interface CameraPreviewProps {
  /** フィルターID */
  filterId?: FilterId | null;
  /** カメラ方向 */
  facingMode?: 'user' | 'environment';
  /** エラーハンドラ */
  onError?: (error: Error) => void;
  /** カメラ起動完了ハンドラ */
  onReady?: (stream: MediaStream) => void;
  /** クラス名 */
  className?: string;
}

export default function CameraPreview({
  filterId = null,
  facingMode = 'environment',
  onError,
  onReady,
  className = '',
}: CameraPreviewProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const rafRef = useRef<number | null>(null);
  const [isReady, setIsReady] = useState(false);

  // カメラ起動
  useEffect(() => {
    let mounted = true;

    const startCamera = async () => {
      try {
        const constraints: MediaStreamConstraints = {
          video: {
            facingMode,
            width: { ideal: 1280 },
            height: { ideal: 720 },
          },
          audio: false,
        };

        const stream = await navigator.mediaDevices.getUserMedia(constraints);

        if (!mounted) {
          stream.getTracks().forEach((track) => track.stop());
          return;
        }

        streamRef.current = stream;

        if (videoRef.current) {
          videoRef.current.srcObject = stream;
          videoRef.current.play();
        }

        setIsReady(true);
        onReady?.(stream);
      } catch (err) {
        console.error('カメラ起動エラー:', err);
        onError?.(err as Error);
      }
    };

    startCamera();

    return () => {
      mounted = false;
      if (rafRef.current) {
        cancelAnimationFrame(rafRef.current);
      }
      if (streamRef.current) {
        streamRef.current.getTracks().forEach((track) => track.stop());
      }
    };
  }, [facingMode, onError, onReady]);

  // フィルター適用ループ
  useEffect(() => {
    if (!isReady || !filterId) return;

    const video = videoRef.current;
    const canvas = canvasRef.current;
    if (!video || !canvas) return;

    const ctx = canvas.getContext('2d', { willReadFrequently: true });
    if (!ctx) return;

    const renderLoop = () => {
      if (video.readyState === video.HAVE_ENOUGH_DATA) {
        canvas.width = video.videoWidth;
        canvas.height = video.videoHeight;

        ctx.drawImage(video, 0, 0, canvas.width, canvas.height);

        if (filterId) {
          const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
          applyFilter(imageData, filterId);
          ctx.putImageData(imageData, 0, 0);
        }
      }

      rafRef.current = requestAnimationFrame(renderLoop);
    };

    renderLoop();

    return () => {
      if (rafRef.current) {
        cancelAnimationFrame(rafRef.current);
      }
    };
  }, [isReady, filterId]);

  return (
    <div className={`${styles.container} ${className}`}>
      {/* フィルターなしの場合はvideoを直接表示 */}
      <video
        ref={videoRef}
        className={styles.video}
        style={{ display: filterId ? 'none' : 'block' }}
        playsInline
        muted
      />

      {/* フィルターありの場合はcanvasを表示 */}
      {filterId && (
        <canvas
          ref={canvasRef}
          className={styles.canvas}
        />
      )}

      {!isReady && (
        <div className={styles.loading}>
          <p>カメラ起動中...</p>
        </div>
      )}
    </div>
  );
}
