/**
 * S9: 2025年 - メッセージ刻み
 *
 * - 最大120文字のメッセージ入力
 * - S2の5箇所から刻む場所を選択
 * - 匿名 + セッションIDのみで保存
 * - 保存したメッセージは後続プレイヤーに提示される
 */

'use client';

import { useState } from 'react';
import styles from './page.module.css';
import Button from '@/components/ui/Button';
import { useGameFlow } from '@/hooks/useGameFlow';
import { saveMessage } from '@/lib/game/api';
import {
  FAVORITE_PLACES,
  MAX_MESSAGE_LENGTH,
} from '@/lib/game/constants';

export default function S9Page() {
  const { session, transitionTo } = useGameFlow();

  const [step, setStep] = useState<'intro' | 'input' | 'select' | 'complete'>(
    'intro'
  );
  const [message, setMessage] = useState('');
  const [selectedPlaceId, setSelectedPlaceId] = useState<string | null>(null);
  const [isSaving, setIsSaving] = useState(false);

  // 入力画面へ
  const handleStartInput = () => {
    setStep('input');
  };

  // メッセージ確定
  const handleConfirmMessage = () => {
    if (message.trim().length < 5) {
      alert('メッセージは5文字以上で入力してください');
      return;
    }

    setStep('select');
  };

  // 場所選択
  const handleSelectPlace = (placeId: string) => {
    setSelectedPlaceId(placeId);
  };

  // メッセージ保存
  const handleSaveMessage = async () => {
    if (!selectedPlaceId || !session) return;

    setIsSaving(true);

    try {
      const success = await saveMessage(session.sessionId, {
        message_text: message.trim(),
        place_id: selectedPlaceId,
      });

      if (success) {
        setStep('complete');
      } else {
        alert('メッセージの保存に失敗しました。もう一度お試しください。');
      }
    } catch (error) {
      console.error('Failed to save message:', error);
      alert('エラーが発生しました。');
    } finally {
      setIsSaving(false);
    }
  };

  const selectedPlace = FAVORITE_PLACES.find((p) => p.id === selectedPlaceId);

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        {/* 年代表示 */}
        <div className={styles.yearBadge}>2025年</div>

        {/* 導入 */}
        {step === 'intro' && (
          <div className={styles.introBox}>
            <h1 className={styles.introTitle}>あなたの痕跡を刻む</h1>

            <div className={styles.introText}>
              <p>
                あなたは、存在証明書を手に入れ、
                <br />
                自分自身の願いを思い出しました。
              </p>
              <p>
                今、あなたは2025年に戻ってきました。
                <br />
                <br />
                この建物に、あなたの痕跡を刻んでください。
                <br />
                未来の誰かが、それを見つけるかもしれません。
              </p>
            </div>

            <Button onClick={handleStartInput} variant="primary" size="large" fullWidth>
              メッセージを刻む
            </Button>
          </div>
        )}

        {/* メッセージ入力 */}
        {step === 'input' && (
          <div className={styles.inputBox}>
            <h2 className={styles.inputTitle}>
              未来の誰かへ、メッセージを残してください
            </h2>

            <textarea
              className={styles.textarea}
              placeholder="あなたの想い、願い、誰かへの言葉..."
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              rows={6}
              maxLength={MAX_MESSAGE_LENGTH}
            />

            <div className={styles.characterCount}>
              {message.length} / {MAX_MESSAGE_LENGTH}
            </div>

            <div className={styles.hint}>
              ※ このメッセージは匿名で保存されます
            </div>

            <Button
              onClick={handleConfirmMessage}
              variant="primary"
              size="large"
              fullWidth
              disabled={message.trim().length < 5}
            >
              次へ
            </Button>
          </div>
        )}

        {/* 場所選択 */}
        {step === 'select' && (
          <div className={styles.selectBox}>
            <h2 className={styles.selectTitle}>
              どの場所にメッセージを刻みますか？
            </h2>

            <div className={styles.messagePreview}>
              <p>{message}</p>
            </div>

            <div className={styles.placesList}>
              {FAVORITE_PLACES.map((place) => (
                <div
                  key={place.id}
                  className={`${styles.placeItem} ${
                    selectedPlaceId === place.id ? styles.selected : ''
                  }`}
                  onClick={() => handleSelectPlace(place.id)}
                >
                  <div className={styles.placeIcon}>📍</div>
                  <div className={styles.placeName}>{place.name}</div>
                </div>
              ))}
            </div>

            <Button
              onClick={handleSaveMessage}
              variant="primary"
              size="large"
              fullWidth
              disabled={!selectedPlaceId}
              isLoading={isSaving}
            >
              {selectedPlace?.name}に刻む
            </Button>
          </div>
        )}

        {/* 完了 */}
        {step === 'complete' && (
          <div className={styles.completeBox}>
            <h1 className={styles.completeTitle}>ありがとうございました</h1>

            <div className={styles.completeText}>
              <p>
                あなたのメッセージは、
                <br />
                {selectedPlace?.name}に刻まれました。
              </p>
              <p>
                いつか、誰かがそれを見つけ、
                <br />
                あなたが確かにここで生きた証を知るでしょう。
              </p>
              <p className={styles.thanks}>
                この建物の歴史に、
                <br />
                あなたも加わりました。
              </p>
            </div>

            <Button
              onClick={() => (window.location.href = '/')}
              variant="secondary"
              size="large"
              fullWidth
            >
              終了
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
