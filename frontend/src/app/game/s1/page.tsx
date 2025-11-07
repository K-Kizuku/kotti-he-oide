/**
 * S1: 1942年パート - 担当者との初回会話
 *
 * - VOICEVOX（青山龍星(しっとり)）による音声生成
 * - プレイヤー情報の入力（来館方法、普段の活動など）
 * - カメラプレビュー + 暗いオーバーレイ
 */

'use client';

import { useState, useEffect, useMemo } from 'react';
import styles from './page.module.css';
import Button from '@/components/ui/Button';
import VoicePlayer from '@/components/audio/VoicePlayer';
import CameraPreview from '@/components/camera/CameraPreview';
import { useGameFlow } from '@/hooks/useGameFlow';

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
  const { transitionTo } = useGameFlow();

  const [step, setStep] = useState<ConversationStep>('welcome');
  const [data, setData] = useState<ConversationData>({
    visitMethod: '',
    activities: '',
  });
  const [currentVoiceUrl, setCurrentVoiceUrl] = useState<string | null>(null);
  const [currentSubtitle, setCurrentSubtitle] = useState('');
  const [hasUserInteracted, setHasUserInteracted] = useState(false);

  // 会話テキスト（定数）
  const dialogues: Record<ConversationStep, string> = useMemo(
    () => ({
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
    }),
    []
  );

  // ステップと音声ファイルのマッピング（定数）
  const voiceFiles: Partial<Record<ConversationStep, string>> = useMemo(
    () => ({
      welcome: '/audio/voice/welcome.wav',
      ask_visit_method: '/audio/voice/ask_visit_method.wav',
      ask_activities: '/audio/voice/ask_activities.wav',
      closing: '/audio/voice/closing.wav',
    }),
    []
  );

  // ステップ変更時に音声を設定
  useEffect(() => {
    const voiceFile = voiceFiles[step];
    const dialogue = dialogues[step];

    // welcome ステップかつユーザー未操作の場合は音声を設定しない
    if (step === 'welcome' && !hasUserInteracted) {
      setCurrentVoiceUrl(null);
      setCurrentSubtitle('');
      return;
    }

    if (voiceFile && dialogue) {
      setCurrentVoiceUrl(voiceFile);
      setCurrentSubtitle(dialogue);
    } else {
      setCurrentVoiceUrl(null);
      setCurrentSubtitle('');
    }
  }, [step, voiceFiles, dialogues, hasUserInteracted]);

  // 会話開始ハンドラ（welcomeステップ用）
  const handleStartConversation = () => {
    setHasUserInteracted(true);
    // 音声を設定（次のuseEffectで自動的に再生される）
    const voiceFile = voiceFiles['welcome'];
    const dialogue = dialogues['welcome'];
    if (voiceFile && dialogue) {
      setCurrentVoiceUrl(voiceFile);
      setCurrentSubtitle(dialogue);
    }
  };

  // 音声終了ハンドラ
  const handleVoiceEnded = () => {
    // 各ステップの音声終了後の自動遷移
    if (step === 'welcome') {
      // welcome音声終了 → 来館方法の質問へ
      setStep('ask_visit_method');
    } else if (step === 'ask_visit_method') {
      // 来館方法の質問終了 → 入力画面へ
      setStep('input_visit_method');
    } else if (step === 'ask_activities') {
      // 活動内容の質問終了 → 入力画面へ
      setStep('input_activities');
    } else if (step === 'closing') {
      // 締めの挨拶終了 → 次へ進むボタンを表示
      setStep('complete');
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

        {/* welcomeステップ：開始ボタン */}
        {step === 'welcome' && !hasUserInteracted && (
          <div className={styles.welcomeBox}>
            <h2 className={styles.welcomeTitle}>1942年 赤煉瓦文化館</h2>
            <p className={styles.welcomeText}>
              担当者があなたを待っています。
              <br />
              会話を開始してください。
            </p>
            <Button
              onClick={handleStartConversation}
              variant="primary"
              size="large"
              fullWidth
            >
              会話を開始する
            </Button>
          </div>
        )}

        {/* 音声プレイヤー（字幕） */}
        {currentVoiceUrl && !isInputStep && (step !== 'welcome' || hasUserInteracted) && (
          <VoicePlayer
            voiceUrl={currentVoiceUrl}
            text={currentSubtitle}
            onEnded={handleVoiceEnded}
            autoPlay
          />
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
        {step === 'complete' && (
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
