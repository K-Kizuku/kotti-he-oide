/**
 * useCamera - カメラ管理フック
 *
 * カメラアクセス、権限管理、ストリーム制御
 */

import { useState, useEffect, useCallback, useRef } from 'react';

interface UseCameraOptions {
  /** カメラの向き */
  facingMode?: 'user' | 'environment';
  /** 自動開始するか */
  autoStart?: boolean;
  /** エラーハンドラ */
  onError?: (error: Error) => void;
}

interface UseCameraReturn {
  /** MediaStreamオブジェクト */
  stream: MediaStream | null;
  /** カメラが起動中か */
  isActive: boolean;
  /** カメラ起動中か（初回読み込み） */
  isLoading: boolean;
  /** エラー */
  error: Error | null;
  /** カメラを開始 */
  startCamera: () => Promise<void>;
  /** カメラを停止 */
  stopCamera: () => void;
  /** カメラの向きを切り替え */
  switchCamera: () => Promise<void>;
}

export function useCamera({
  facingMode = 'environment',
  autoStart = false,
  onError,
}: UseCameraOptions = {}): UseCameraReturn {
  const [stream, setStream] = useState<MediaStream | null>(null);
  const [isActive, setIsActive] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [currentFacingMode, setCurrentFacingMode] = useState(facingMode);
  const streamRef = useRef<MediaStream | null>(null);

  const stopCamera = useCallback(() => {
    if (streamRef.current) {
      streamRef.current.getTracks().forEach((track) => {
        track.stop();
      });
      streamRef.current = null;
      setStream(null);
      setIsActive(false);
    }
  }, []);

  const startCamera = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);

      // 既存のストリームを停止
      stopCamera();

      const mediaStream = await navigator.mediaDevices.getUserMedia({
        video: {
          facingMode: currentFacingMode,
          width: { ideal: 1920 },
          height: { ideal: 1080 },
        },
        audio: false,
      });

      streamRef.current = mediaStream;
      setStream(mediaStream);
      setIsActive(true);
    } catch (err) {
      const error =
        err instanceof Error ? err : new Error('Failed to access camera');
      setError(error);
      onError?.(error);
      console.error('Camera access error:', error);
    } finally {
      setIsLoading(false);
    }
  }, [currentFacingMode, stopCamera, onError]);

  const switchCamera = useCallback(async () => {
    const newFacingMode = currentFacingMode === 'user' ? 'environment' : 'user';
    setCurrentFacingMode(newFacingMode);

    if (isActive) {
      stopCamera();
      // 少し待ってから新しいカメラを起動
      setTimeout(() => {
        startCamera();
      }, 100);
    }
  }, [currentFacingMode, isActive, stopCamera, startCamera]);

  // 自動開始
  useEffect(() => {
    if (autoStart) {
      startCamera();
    }

    // クリーンアップ
    return () => {
      stopCamera();
    };
  }, [autoStart, startCamera, stopCamera]);

  return {
    stream,
    isActive,
    isLoading,
    error,
    startCamera,
    stopCamera,
    switchCamera,
  };
}
