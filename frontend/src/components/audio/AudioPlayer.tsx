/**
 * AudioPlayer - BGM/SE再生コンポーネント
 *
 * 背景音楽や効果音を再生
 * ループ、音量調整などをサポート
 */

'use client';

import { useEffect, useRef } from 'react';

interface AudioPlayerProps {
  /** 音声ファイルURL */
  src: string;
  /** 自動再生するか */
  autoPlay?: boolean;
  /** ループ再生するか */
  loop?: boolean;
  /** 音量（0-1） */
  volume?: number;
  /** 再生開始ハンドラ */
  onPlay?: () => void;
  /** 再生終了ハンドラ */
  onEnded?: () => void;
  /** エラーハンドラ */
  onError?: (error: Error) => void;
}

export default function AudioPlayer({
  src,
  autoPlay = false,
  loop = false,
  volume = 0.7,
  onPlay,
  onEnded,
  onError,
}: AudioPlayerProps) {
  const audioRef = useRef<HTMLAudioElement>(null);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    audio.volume = Math.max(0, Math.min(1, volume));

    if (autoPlay) {
      audio.play().catch((err) => {
        console.error('音声再生エラー:', err);
        onError?.(err);
      });
    }
  }, [volume, autoPlay, onError]);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const handlePlay = () => onPlay?.();
    const handleEnded = () => onEnded?.();
    const handleError = () => onError?.(new Error('Audio playback error'));

    audio.addEventListener('play', handlePlay);
    audio.addEventListener('ended', handleEnded);
    audio.addEventListener('error', handleError);

    return () => {
      audio.removeEventListener('play', handlePlay);
      audio.removeEventListener('ended', handleEnded);
      audio.removeEventListener('error', handleError);
    };
  }, [onPlay, onEnded, onError]);

  return (
    <audio
      ref={audioRef}
      src={src}
      loop={loop}
      preload="auto"
      style={{ display: 'none' }}
    />
  );
}
