/**
 * useAudio - 音声再生管理フック
 *
 * BGM/SE/VOICEVOXの再生を管理
 */

import { useState, useEffect, useRef, useCallback } from 'react';

interface UseAudioOptions {
  /** 音声URL */
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

export function useAudio({
  src,
  autoPlay = false,
  loop = false,
  volume = 0.7,
  onPlay,
  onEnded,
  onError,
}: UseAudioOptions) {
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  // Audio要素の初期化
  useEffect(() => {
    const audio = new Audio(src);
    audio.loop = loop;
    audio.volume = Math.max(0, Math.min(1, volume));

    audioRef.current = audio;

    const handlePlay = () => {
      setIsPlaying(true);
      onPlay?.();
    };

    const handleEnded = () => {
      setIsPlaying(false);
      onEnded?.();
    };

    const handleError = () => {
      const err = new Error('Audio playback error');
      setError(err);
      setIsLoading(false);
      onError?.(err);
    };

    const handleCanPlay = () => {
      setIsLoading(false);
    };

    audio.addEventListener('play', handlePlay);
    audio.addEventListener('ended', handleEnded);
    audio.addEventListener('error', handleError);
    audio.addEventListener('canplaythrough', handleCanPlay);

    if (autoPlay) {
      audio.play().catch((err) => {
        setError(err);
        onError?.(err);
      });
    }

    return () => {
      audio.removeEventListener('play', handlePlay);
      audio.removeEventListener('ended', handleEnded);
      audio.removeEventListener('error', handleError);
      audio.removeEventListener('canplaythrough', handleCanPlay);
      audio.pause();
      audio.src = '';
    };
  }, [src, loop, volume, autoPlay, onPlay, onEnded, onError]);

  const play = useCallback(() => {
    if (audioRef.current && !isPlaying) {
      audioRef.current.play().catch((err) => {
        setError(err);
        onError?.(err);
      });
    }
  }, [isPlaying, onError]);

  const pause = useCallback(() => {
    if (audioRef.current && isPlaying) {
      audioRef.current.pause();
    }
  }, [isPlaying]);

  const stop = useCallback(() => {
    if (audioRef.current) {
      audioRef.current.pause();
      audioRef.current.currentTime = 0;
    }
  }, []);

  const setVolume = useCallback((newVolume: number) => {
    if (audioRef.current) {
      audioRef.current.volume = Math.max(0, Math.min(1, newVolume));
    }
  }, []);

  return {
    isPlaying,
    isLoading,
    error,
    play,
    pause,
    stop,
    setVolume,
  };
}
