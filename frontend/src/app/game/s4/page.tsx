/**
 * S4: 診査室 - 内省質問パート
 *
 * - 固定10問の質問（小学生時代の夢中だったこと、尊敬した人、人生の願望など）
 * - "なし"や"特にない"などの無回答はバリデーションで弾く
 * - 回答は逐次サーバーに保存（セッション再開対応）
 * - 必須質問：「人生の最期に達成したいこと」「名前」
 */

'use client';

import { useState, useEffect } from 'react';
import styles from './page.module.css';
import Button from '@/components/ui/Button';
import ProgressBar from '@/components/ui/ProgressBar';
import { useGameFlow } from '@/hooks/useGameFlow';
import { saveAnswer, getAnswers } from '@/lib/game/api';
import {
  INTROSPECTION_QUESTIONS,
  isInvalidAnswer,
} from '@/lib/game/constants';

export default function S4Page() {
  const { session, transitionTo } = useGameFlow();

  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [answers, setAnswers] = useState<Record<string, string>>({});
  const [currentAnswer, setCurrentAnswer] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState('');

  const currentQuestion = INTROSPECTION_QUESTIONS[currentQuestionIndex];
  const progress =
    ((currentQuestionIndex + 1) / INTROSPECTION_QUESTIONS.length) * 100;

  // 保存済み回答を読み込む
  useEffect(() => {
    if (!session) return;

    const loadSavedAnswers = async () => {
      const savedAnswers = await getAnswers(session.sessionId);
      const answersMap: Record<string, string> = {};

      savedAnswers.forEach((saved) => {
        answersMap[saved.question_id] = saved.answer_text;
      });

      setAnswers(answersMap);

      // 最後に回答した質問の次から開始
      const lastAnsweredIndex = INTROSPECTION_QUESTIONS.findIndex(
        (q) => !answersMap[q.id]
      );
      if (lastAnsweredIndex !== -1) {
        setCurrentQuestionIndex(lastAnsweredIndex);
      } else if (Object.keys(answersMap).length === INTROSPECTION_QUESTIONS.length) {
        // 全て回答済みの場合は最後の質問へ
        setCurrentQuestionIndex(INTROSPECTION_QUESTIONS.length - 1);
        setCurrentAnswer(answersMap[INTROSPECTION_QUESTIONS[INTROSPECTION_QUESTIONS.length - 1].id] || '');
      }
    };

    loadSavedAnswers();
  }, [session]);

  // 現在の質問の回答を復元
  useEffect(() => {
    if (answers[currentQuestion.id]) {
      setCurrentAnswer(answers[currentQuestion.id]);
    } else {
      setCurrentAnswer('');
    }
  }, [currentQuestionIndex, currentQuestion.id, answers]);

  // バリデーション
  const validateAnswer = (answer: string): string | null => {
    const trimmed = answer.trim();

    if (!trimmed) {
      return '回答を入力してください';
    }

    if (currentQuestion.validation?.minLength && trimmed.length < currentQuestion.validation.minLength) {
      return `${currentQuestion.validation.minLength}文字以上で入力してください`;
    }

    if (currentQuestion.validation?.maxLength && trimmed.length > currentQuestion.validation.maxLength) {
      return `${currentQuestion.validation.maxLength}文字以内で入力してください`;
    }

    if (isInvalidAnswer(trimmed)) {
      return '「なし」「特になし」などの回答は受け付けられません。具体的に答えてください。';
    }

    return null;
  };

  // 回答を保存して次へ
  const handleNext = async () => {
    setError('');

    const validationError = validateAnswer(currentAnswer);
    if (validationError) {
      setError(validationError);
      return;
    }

    if (!session) {
      setError('セッションが見つかりません');
      return;
    }

    setIsSaving(true);

    try {
      // サーバーに保存
      const success = await saveAnswer(session.sessionId, {
        question_id: currentQuestion.id,
        answer_text: currentAnswer.trim(),
      });

      if (!success) {
        throw new Error('Failed to save answer');
      }

      // ローカル状態を更新
      setAnswers((prev) => ({
        ...prev,
        [currentQuestion.id]: currentAnswer.trim(),
      }));

      // 次の質問へ
      if (currentQuestionIndex < INTROSPECTION_QUESTIONS.length - 1) {
        setCurrentQuestionIndex((prev) => prev + 1);
        setCurrentAnswer('');
      } else {
        // 全問回答済み
        await transitionTo('s5');
      }
    } catch (err) {
      console.error('Failed to save answer:', err);
      setError('回答の保存に失敗しました。もう一度お試しください。');
    } finally {
      setIsSaving(false);
    }
  };

  // 前の質問に戻る
  const handlePrevious = () => {
    if (currentQuestionIndex > 0) {
      setCurrentQuestionIndex((prev) => prev - 1);
      setError('');
    }
  };

  const isLastQuestion = currentQuestionIndex === INTROSPECTION_QUESTIONS.length - 1;

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        {/* ヘッダー */}
        <header className={styles.header}>
          <h1 className={styles.title}>診査室</h1>
          <p className={styles.subtitle}>
            あなた自身について、お聞かせください
          </p>
        </header>

        {/* 進捗バー */}
        <ProgressBar
          value={progress}
          max={100}
          label={`${currentQuestionIndex + 1} / ${INTROSPECTION_QUESTIONS.length}`}
        />

        {/* 質問カード */}
        <div className={styles.questionCard}>
          <div className={styles.questionNumber}>
            質問 {currentQuestionIndex + 1}
          </div>

          <div className={styles.questionText}>
            <p>{currentQuestion.text}</p>
          </div>

          <textarea
            className={styles.textarea}
            placeholder={currentQuestion.placeholder || '回答を入力してください...'}
            value={currentAnswer}
            onChange={(e) => setCurrentAnswer(e.target.value)}
            rows={5}
            maxLength={currentQuestion.validation?.maxLength || 200}
            disabled={isSaving}
          />

          <div className={styles.characterCount}>
            {currentAnswer.length} / {currentQuestion.validation?.maxLength || 200}
          </div>

          {/* エラーメッセージ */}
          {error && (
            <div className={styles.errorBox}>
              <p>{error}</p>
            </div>
          )}
        </div>

        {/* ナビゲーション */}
        <div className={styles.navigation}>
          {currentQuestionIndex > 0 && (
            <Button
              onClick={handlePrevious}
              variant="secondary"
              size="medium"
              disabled={isSaving}
            >
              ← 前の質問
            </Button>
          )}

          <div className={styles.spacer} />

          <Button
            onClick={handleNext}
            variant="primary"
            size="large"
            disabled={isSaving || !currentAnswer.trim()}
            isLoading={isSaving}
          >
            {isLastQuestion ? '完了' : '次へ →'}
          </Button>
        </div>
      </div>
    </div>
  );
}
