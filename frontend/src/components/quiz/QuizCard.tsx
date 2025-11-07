/**
 * QuizCard - クイズカードコンポーネント
 *
 * S6で使用される4択クイズ
 */

'use client';

import { useState } from 'react';
import styles from './QuizCard.module.css';
import Button from '../ui/Button';

export interface QuizOption {
  id: string;
  label: string;
}

interface QuizCardProps {
  /** 質問文 */
  question: string;
  /** 選択肢（4つ） */
  options: QuizOption[];
  /** 正解の選択肢ID */
  correctAnswerId: string;
  /** 回答ハンドラ */
  onAnswer: (isCorrect: boolean, selectedId: string) => void;
  /** 場所名（オプション） */
  placeName?: string;
}

export default function QuizCard({
  question,
  options,
  correctAnswerId,
  onAnswer,
  placeName,
}: QuizCardProps) {
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [hasAnswered, setHasAnswered] = useState(false);

  const handleSelect = (optionId: string) => {
    if (hasAnswered) return;
    setSelectedId(optionId);
  };

  const handleSubmit = () => {
    if (!selectedId || hasAnswered) return;

    setHasAnswered(true);
    const isCorrect = selectedId === correctAnswerId;
    onAnswer(isCorrect, selectedId);
  };

  return (
    <div className={styles.card}>
      {placeName && (
        <div className={styles.placeName}>
          {placeName}
        </div>
      )}

      <div className={styles.question}>
        <p>{question}</p>
      </div>

      <div className={styles.options}>
        {options.map((option) => {
          const isSelected = selectedId === option.id;
          const isCorrect = hasAnswered && option.id === correctAnswerId;
          const isWrong = hasAnswered && isSelected && option.id !== correctAnswerId;

          return (
            <button
              key={option.id}
              className={`${styles.option} ${isSelected ? styles.selected : ''} ${isCorrect ? styles.correct : ''} ${isWrong ? styles.wrong : ''}`}
              onClick={() => handleSelect(option.id)}
              disabled={hasAnswered}
            >
              {option.label}
            </button>
          );
        })}
      </div>

      <div className={styles.footer}>
        <Button
          onClick={handleSubmit}
          disabled={!selectedId || hasAnswered}
          variant="primary"
          size="large"
          fullWidth
        >
          回答する
        </Button>
      </div>
    </div>
  );
}
