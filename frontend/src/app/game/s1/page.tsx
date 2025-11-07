/**
 * S1: 1942年パート - 担当者との初回会話
 *
 * - VOICEVOX（青山龍星(しっとり)）による音声生成
 * - プレイヤー情報の入力（来館方法、普段の活動など）
 * - カメラプレビュー + 暗いオーバーレイ
 */

'use client';

import { useState, useEffect } from 'react';
import styles from './page.module.css';
import Button from '@/components/ui/Button';
import VoicePlayer from '@/components/audio/VoicePlayer';
import CameraPreview from '@/components/camera/CameraPreview';
import { useGameFlow } from '@/hooks/useGameFlow';
import { generateVoice } from '@/lib/game/api';

type ConversationStep =
  | 'welcome'
  | 'ask_visit_method'
  | 'input_visit_method'
  | 'ask_activities'
  | 'input_activities'
  | 'closing'
  | 'complete';

interface ConversationData {
  visitMethod: string;
  activities: string;
}

export default function S1Page() {
  const { session, transitionTo } = useGameFlow();

  const [step, setStep] = useState<ConversationStep>('welcome');
  const [data, setData] = useState<ConversationData>({
    visitMethod: '',
    activities: '',
  });
  const [currentVoiceUrl, setCurrentVoiceUrl] = useState<string | null>(null);
  const [isLoadingVoice, setIsLoadingVoice] = useState(false);
  const [currentSubtitle, setCurrentSubtitle] = useState('');

  // 会話テキスト
  const dialogues: Record<ConversationStep, string> = {
    welcome:
      'ようこそ、赤煉瓦文化館へ。私は本日の担当者でございます。お越しいただき誠にありがとうございます。',
    ask_visit_method:
      '本日はどのようにしてこちらまでいらっしゃいましたか？',
    input_visit_method: '',
    ask_activities: 'なるほど。では、普段はどのようなことをされていますか？',
    input_activities: '',
    closing:
      'ありがとうございます。それでは、まずこの建物でお気に入りの場所を見つけていただきたいと思います。',
    complete: '',
  };

  // 音声生成
  const playVoice = async (text: string) => {
    if (!text) return;

    setIsLoadingVoice(true);
    setCurrentSubtitle(text);

    try {
      const audioUrl = await generateVoice({ text });
      if (audioUrl) {
        setCurrentVoiceUrl(audioUrl);
      }
    } catch (error) {
      console.error('Failed to generate voice:', error);
    } finally {
      setIsLoadingVoice(false);
    }
  };

  // ステップ変更時に音声を再生
  useEffect(() => {
    const dialogue = dialogues[step];
    if (dialogue) {
      playVoice(dialogue);
    }
  }, [step]);

  // 音声終了ハンドラ
  const handleVoiceEnded = () => {
    // 入力ステップに自動遷移
    if (step === 'ask_visit_method') {
      setStep('input_visit_method');
    } else if (step === 'ask_activities') {
      setStep('input_activities');
    }
  };

  // 来館方法の送信
  const handleVisitMethodSubmit = () => {
    if (data.visitMethod.trim().length < 2) {
      alert('来館方法を入力してください');
      return;
    }

    setStep('ask_activities');
  };

  // 活動内容の送信
  const handleActivitiesSubmit = () => {
    if (data.activities.trim().length < 2) {
      alert('普段の活動を入力してください');
      return;
    }

    setStep('closing');
  };

  // 会話終了
  const handleComplete = async () => {
    await transitionTo('s2');
  };

  // 入力ステップかどうか
  const isInputStep =
    step === 'input_visit_method' || step === 'input_activities';

  return (
    <div className={styles.container}>
      {/* カメラプレビュー（背景） */}
      <div className={styles.cameraBackground}>
        <CameraPreview facingMode="environment" />
        <div className={styles.darkOverlay} />
      </div>

      {/* コンテンツ */}
      <div className={styles.content}>
        {/* 年代表示 */}
        <div className={styles.yearBadge}>1942年</div>

        {/* 音声プレイヤー（字幕） */}
        {currentVoiceUrl && !isInputStep && (
          <VoicePlayer
            voiceUrl={currentVoiceUrl}
            text={currentSubtitle}
            onEnded={handleVoiceEnded}
            autoPlay
          />
        )}

        {/* ローディング */}
        {isLoadingVoice && (
          <div className={styles.loadingBox}>
            <p>音声を生成しています...</p>
          </div>
        )}

        {/* 入力フォーム */}
        {step === 'input_visit_method' && (
          <div className={styles.inputSection}>
            <div className={styles.questionBox}>
              <p className={styles.question}>{dialogues.ask_visit_method}</p>
            </div>

            <textarea
              className={styles.textarea}
              placeholder="例：徒歩で、電車とバスで、自動車で..."
              value={data.visitMethod}
              onChange={(e) =>
                setData({ ...data, visitMethod: e.target.value })
              }
              rows={3}
              maxLength={100}
            />

            <Button
              onClick={handleVisitMethodSubmit}
              variant="primary"
              size="large"
              fullWidth
              disabled={data.visitMethod.trim().length < 2}
            >
              回答する
            </Button>
          </div>
        )}

        {step === 'input_activities' && (
          <div className={styles.inputSection}>
            <div className={styles.questionBox}>
              <p className={styles.question}>{dialogues.ask_activities}</p>
            </div>

            <textarea
              className={styles.textarea}
              placeholder="例：会社で働いている、学生です、家事をしています..."
              value={data.activities}
              onChange={(e) =>
                setData({ ...data, activities: e.target.value })
              }
              rows={3}
              maxLength={100}
            />

            <Button
              onClick={handleActivitiesSubmit}
              variant="primary"
              size="large"
              fullWidth
              disabled={data.activities.trim().length < 2}
            >
              回答する
            </Button>
          </div>
        )}

        {/* 完了ボタン */}
        {step === 'closing' && !isLoadingVoice && (
          <div className={styles.nextSection}>
            <Button
              onClick={handleComplete}
              variant="primary"
              size="large"
              fullWidth
            >
              次へ進む
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
