/**
 * S7: 2002年パート - メッセージ受け取り
 *
 * - S4で答えた「人生の最期に達成したいこと」を一字一句で再入力
 * - 不一致の場合は再入力を要求
 */

'use client';

import { useState, useEffect } from 'react';
import styles from './page.module.css';
import Button from '@/components/ui/Button';
import { useGameFlow } from '@/hooks/useGameFlow';
import { getAnswers, type SavedAnswer } from '@/lib/game/api';

export default function S7Page() {
  const { session, transitionTo } = useGameFlow();

  const [lifeGoalAnswer, setLifeGoalAnswer] = useState('');
  const [userInput, setUserInput] = useState('');
  const [error, setError] = useState('');
  const [step, setStep] = useState<'intro' | 'input' | 'success'>('intro');

  // S4の回答を読み込む
  useEffect(() => {
    if (!session) return;

    const loadAnswers = async () => {
      const answers = await getAnswers(session.sessionId);
      const lifeGoal = answers.find((a) => a.question_id === 'q9_life_goal');

      if (lifeGoal) {
        setLifeGoalAnswer(lifeGoal.answer_text);
      }
    };

    loadAnswers();
  }, [session]);

  // 導入から入力画面へ
  const handleStartInput = () => {
    setStep('input');
  };

  // 回答検証
  const handleSubmit = () => {
    setError('');

    const trimmedInput = userInput.trim();
    const trimmedAnswer = lifeGoalAnswer.trim();

    if (trimmedInput === trimmedAnswer) {
      // 完全一致：成功
      setStep('success');
      setTimeout(async () => {
        await transitionTo('s8');
      }, 3000);
    } else {
      // 不一致：エラー
      setError(
        'あなたの答えと一致しません。\n一字一句、正確に思い出してください。'
      );
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        {/* 年代表示 */}
        <div className={styles.yearBadge}>2002年</div>

        {/* 導入 */}
        {step === 'intro' && (
          <div className={styles.introBox}>
            <h1 className={styles.introTitle}>あなたの願いを思い出せ</h1>

            <div className={styles.introText}>
              <p>
                存在証明書を手に入れたあなたは、
                <br />
                再び生きる権利を得ました。
              </p>
              <p>
                しかし、それには条件があります。
                <br />
                <br />
                あなた自身が語った、
                <br />
                人生の最期に達成したいことを、
                <br />
                一字一句正確に思い出してください。
              </p>
            </div>

            <Button onClick={handleStartInput} variant="primary" size="large" fullWidth>
              思い出す
            </Button>
          </div>
        )}

        {/* 入力画面 */}
        {step === 'input' && (
          <div className={styles.inputBox}>
            <h2 className={styles.questionTitle}>
              人生の最期に達成したいことは何ですか？
            </h2>

            <textarea
              className={styles.textarea}
              placeholder="あなたが答えたことを、正確に入力してください..."
              value={userInput}
              onChange={(e) => setUserInput(e.target.value)}
              rows={4}
              maxLength={100}
            />

            <div className={styles.hint}>
              ※ 一字一句、正確に入力してください
            </div>

            {error && (
              <div className={styles.errorBox}>
                <p>{error}</p>
              </div>
            )}

            <Button
              onClick={handleSubmit}
              variant="primary"
              size="large"
              fullWidth
              disabled={!userInput.trim()}
            >
              確認する
            </Button>
          </div>
        )}

        {/* 成功 */}
        {step === 'success' && (
          <div className={styles.successBox}>
            <h2 className={styles.successTitle}>思い出しました</h2>

            <div className={styles.successMessage}>
              <p className={styles.goalText}>「{lifeGoalAnswer}」</p>

              <p className={styles.successDescription}>
                その願いを、決して忘れないでください。
                <br />
                <br />
                それがあなたの、生きる証です。
              </p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
