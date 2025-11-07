/**
 * S8: メインホール3分探索
 *
 * - マップ画像から該当箇所を探す
 * - 3分カウントダウン
 * - 成功時に過去プレイヤーのメッセージ一覧を表示
 */

'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import styles from './page.module.css';
import Button from '@/components/ui/Button';
import Timer from '@/components/ui/Timer';
import { useTimer } from '@/hooks/useTimer';
import { getMessages, type PlayerMessage } from '@/lib/game/api';
import {
  TIMER_DURATIONS,
  TIMER_WARNING_THRESHOLDS,
} from '@/lib/game/constants';

export default function S8Page() {
  const router = useRouter();
  const [step, setStep] = useState<'search' | 'messages'>('search');
  const [messages, setMessages] = useState<PlayerMessage[]>([]);

  const timer = useTimer({
    initialSeconds: TIMER_DURATIONS.S8_EXPLORATION,
    autoStart: true,
    warningThreshold: TIMER_WARNING_THRESHOLDS.S8,
    onComplete: () => {
      if (step !== 'messages') {
        router.push('/game/gameover');
      }
    },
  });

  // 発見成功
  const handleFound = async () => {
    // タイマーを停止
    timer.pause();

    // 過去プレイヤーのメッセージを取得
    const fetchedMessages = await getMessages();
    setMessages(fetchedMessages);
    setStep('messages');
  };

  // 次へ進む
  const handleNext = () => {
    router.push('/game/s9');
  };

  return (
    <div className={styles.container}>
      {/* タイマー（固定表示） */}
      {step === 'search' && (
        <Timer
          initialSeconds={timer.seconds}
          large
          fixed
          warningThreshold={TIMER_WARNING_THRESHOLDS.S8}
        />
      )}

      <div className={styles.content}>
        {/* 探索パート */}
        {step === 'search' && (
          <>
            <header className={styles.header}>
              <h1 className={styles.title}>メインホールを探索</h1>
              <p className={styles.subtitle}>
                過去のプレイヤーが残した痕跡を見つけてください
              </p>
            </header>

            <div className={styles.mapContainer}>
              {/* マップ画像（仮） */}
              <div className={styles.mapPlaceholder}>
                <p>館内マップ</p>
                <p className={styles.mapHint}>
                  ※ 実際のゲームでは館内マップ画像を表示
                </p>
              </div>

              <Button
                onClick={handleFound}
                variant="primary"
                size="large"
                fullWidth
              >
                痕跡を発見した
              </Button>
            </div>
          </>
        )}

        {/* メッセージ一覧パート */}
        {step === 'messages' && (
          <>
            <header className={styles.header}>
              <h1 className={styles.title}>過去の声</h1>
              <p className={styles.subtitle}>
                ここで生きた人々の、想いの痕跡
              </p>
            </header>

            <div className={styles.messagesContainer}>
              {messages.length === 0 ? (
                <div className={styles.noMessages}>
                  <p>まだ誰もメッセージを残していません</p>
                  <p>あなたが最初の一人になります</p>
                </div>
              ) : (
                <div className={styles.messagesList}>
                  {messages.map((message) => (
                    <div key={message.message_id} className={styles.messageCard}>
                      <p className={styles.messageText}>
                        {message.message_text}
                      </p>
                      <span className={styles.messageDate}>
                        {new Date(message.created_at).toLocaleDateString('ja-JP', {
                          year: 'numeric',
                          month: 'long',
                          day: 'numeric',
                        })}
                      </span>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <Button
              onClick={handleNext}
              variant="primary"
              size="large"
              fullWidth
            >
              次へ進む
            </Button>
          </>
        )}
      </div>
    </div>
  );
}
