/**
 * VoicePlayer - VOICEVOX音声再生コンポーネント
 *
 * S1で使用
 * 担当者のセリフを再生
 */

'use client';

import { useEffect, useRef, useState } from 'react';
import styles from './VoicePlayer.module.css';

interface VoicePlayerProps {
  /** 音声ファイルURL */
  voiceUrl: string;
  /** 表示するテキスト（字幕） */
  text?: string;
  /** 自動再生するか */
  autoPlay?: boolean;
  /** 再生終了ハンドラ */
  onEnded?: () => void;
  /** 字幕を表示するか */
  showSubtitle?: boolean;
}

export default function VoicePlayer({
  voiceUrl,
  text,
  autoPlay = true,
  onEnded,
  showSubtitle = true,
}: VoicePlayerProps) {
  const audioRef = useRef<HTMLAudioElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    if (autoPlay) {
      audio.play().catch((err) => {
        console.error('音声再生エラー:', err);
      });
    }

    const handlePlay = () => setIsPlaying(true);
    const handleEnded = () => {
      setIsPlaying(false);
      onEnded?.();
    };

    audio.addEventListener('play', handlePlay);
    audio.addEventListener('ended', handleEnded);

    return () => {
      audio.removeEventListener('play', handlePlay);
      audio.removeEventListener('ended', handleEnded);
    };
  }, [voiceUrl, autoPlay, onEnded]);

  return (
    <div className={styles.container}>
      <audio
        ref={audioRef}
        src={voiceUrl}
        preload="auto"
        style={{ display: 'none' }}
      />

      {showSubtitle && text && (
        <div className={`${styles.subtitle} ${isPlaying ? styles.active : ''}`}>
          <p>{text}</p>
        </div>
      )}
    </div>
  );
}
