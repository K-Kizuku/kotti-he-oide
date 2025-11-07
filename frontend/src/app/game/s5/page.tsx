/**
 * S5: 死亡届受理 + 7分制限開始（1972年パート）
 *
 * - 「死亡届が受理されました」の通知
 * - 7分カウントダウン開始
 * - 担当者の存在証明書との交換を要求
 */

'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import styles from './page.module.css';
import Button from '@/components/ui/Button';
import GlitchText from '@/components/horror/GlitchText';
import { useSession } from '@/hooks/useSession';
import { startS6Exploration } from '@/lib/game/api';
import { TIMER_DURATIONS } from '@/lib/game/constants';

export default function S5Page() {
  const router = useRouter();
  const { session } = useSession();
  const [step, setStep] = useState<'notification' | 'explanation' | 'ready'>(
    'notification'
  );

  useEffect(() => {
    // 3秒後に説明へ
    const timer1 = setTimeout(() => {
      setStep('explanation');
    }, 3000);

    // 6秒後に準備完了へ
    const timer2 = setTimeout(() => {
      setStep('ready');
    }, 6000);

    return () => {
      clearTimeout(timer1);
      clearTimeout(timer2);
    };
  }, []);

  const handleStart = async () => {
    if (!session) return;

    // S6探索を開始（サーバー側でタイマー記録）
    await startS6Exploration(session.sessionId);

    router.push('/game/s6');
  };

  const minutes = Math.floor(TIMER_DURATIONS.S6_EXPLORATION / 60);

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        {/* 年代表示 */}
        <div className={styles.yearBadge}>1972年</div>

        {/* 死亡届受理通知 */}
        {step === 'notification' && (
          <div className={styles.notificationBox}>
            <GlitchText intensity="high">
              <span className={styles.notificationText}>
                死亡届が受理されました
              </span>
            </GlitchText>
          </div>
        )}

        {/* 説明 */}
        {step === 'explanation' && (
          <div className={styles.explanationBox}>
            <p className={styles.explanationText}>
              あなたは既に、この世にはいません。
              <br />
              <br />
              しかし、私の存在証明書と引き換えに、
              <br />
              もう一度、生きる機会を差し上げましょう。
            </p>
          </div>
        )}

        {/* 準備完了 */}
        {step === 'ready' && (
          <div className={styles.readyBox}>
            <h2 className={styles.taskTitle}>あなたの課題</h2>

            <div className={styles.taskDescription}>
              <p>
                館内の5つの場所を訪れ、
                <br />
                それぞれで存在証明書のピースを集めてください。
              </p>

              <div className={styles.timeLimit}>
                <span className={styles.timeLimitLabel}>制限時間</span>
                <span className={styles.timeLimitValue}>{minutes}分</span>
              </div>

              <p className={styles.warning}>
                時間内に全てのピースを集められなければ、
                <br />
                あなたの存在は消滅します。
              </p>
            </div>

            <Button
              onClick={handleStart}
              variant="danger"
              size="large"
              fullWidth
            >
              探索を開始する
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
