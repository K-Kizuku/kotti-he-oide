/**
 * S6: 存在証明書探索パート（最重要）
 *
 * - 目標: 5つの場所すべてで1枚ずつピースを取得（合計5枚）
 * - 制限時間: 7分以内
 * - 到達判定: カメラ撮影 or Web選択
 * - クイズシステム: S4の回答を元に生成
 * - 正解でピース取得、不正解は即再挑戦可能（ホラー演出あり）
 */

'use client';

import { useState, useEffect } from 'react';
import styles from './page.module.css';
import Button from '@/components/ui/Button';
import Timer from '@/components/ui/Timer';
import Modal from '@/components/ui/Modal';
import QuizCard from '@/components/quiz/QuizCard';
import CameraCapture from '@/components/camera/CameraCapture';
import JumpScare from '@/components/horror/JumpScare';
import { useGameFlow } from '@/hooks/useGameFlow';
import { useTimer } from '@/hooks/useTimer';
import {
  getQuiz,
  answerQuiz,
  verifyLocation,
  getS6Progress,
  type Quiz,
} from '@/lib/game/api';
import {
  FAVORITE_PLACES,
  TIMER_DURATIONS,
  TIMER_WARNING_THRESHOLDS,
} from '@/lib/game/constants';

type PlaceStatus = {
  verified: boolean; // 場所に到達済み
  answered: boolean; // クイズ回答済み
  correct: boolean; // 正解したか
};

export default function S6Page() {
  const { session, transitionTo } = useGameFlow();

  // 各場所の状態
  const [placeStatuses, setPlaceStatuses] = useState<
    Record<string, PlaceStatus>
  >({});

  // 現在選択中の場所
  const [selectedPlaceId, setSelectedPlaceId] = useState<string | null>(null);

  // モーダル状態
  const [showCamera, setShowCamera] = useState(false);
  const [showQuiz, setShowQuiz] = useState(false);
  const [currentQuiz, setCurrentQuiz] = useState<Quiz | null>(null);

  // ジャンプスケア
  const [showJumpScare, setShowJumpScare] = useState(false);

  // タイマー
  const timer = useTimer({
    initialSeconds: TIMER_DURATIONS.S6_EXPLORATION,
    autoStart: true,
    warningThreshold: TIMER_WARNING_THRESHOLDS.S6,
    onComplete: handleTimerComplete,
  });

  // 進捗を読み込む
  useEffect(() => {
    if (!session) return;

    const loadProgress = async () => {
      const progress = await getS6Progress(session.sessionId);
      const statuses: Record<string, PlaceStatus> = {};

      progress.forEach((p) => {
        statuses[p.place_id] = {
          verified: p.verified,
          answered: p.answered,
          correct: p.correct,
        };
      });

      setPlaceStatuses(statuses);
    };

    loadProgress();
  }, [session]);

  // タイマー完了ハンドラ
  function handleTimerComplete() {
    // 全ピース取得済みなら次へ進める
    const allCompleted = FAVORITE_PLACES.every(
      (place) => placeStatuses[place.id]?.correct
    );

    if (!allCompleted) {
      await transitionTo('gameover');
    }
  }

  // 場所選択
  const handleSelectPlace = (placeId: string) => {
    setSelectedPlaceId(placeId);
  };

  // カメラ撮影モード
  const handleOpenCamera = () => {
    setShowCamera(true);
  };

  // Web選択（手動確認）
  const handleManualVerify = async () => {
    if (!selectedPlaceId || !session) return;

    // 手動確認として記録
    setPlaceStatuses((prev) => ({
      ...prev,
      [selectedPlaceId]: {
        ...prev[selectedPlaceId],
        verified: true,
      },
    }));

    // クイズを取得して表示
    await loadAndShowQuiz(selectedPlaceId);
  };

  // カメラで撮影完了
  const handleCaptureComplete = async (imageBlob: Blob) => {
    if (!selectedPlaceId || !session) return;

    setShowCamera(false);

    // 画像類似度判定
    const result = await verifyLocation(session.sessionId, {
      place_id: selectedPlaceId,
      image: imageBlob,
    });

    if (result.verified) {
      // 到達成功
      setPlaceStatuses((prev) => ({
        ...prev,
        [selectedPlaceId]: {
          ...prev[selectedPlaceId],
          verified: true,
        },
      }));

      // クイズを表示
      await loadAndShowQuiz(selectedPlaceId);
    } else {
      // 到達失敗
      alert(
        `この場所の写真ではないようです（類似度: ${(result.similarity * 100).toFixed(0)}%）\n\nもう一度撮影するか、「この場所にいることにする」を選択してください。`
      );
    }
  };

  // クイズを読み込んで表示
  const loadAndShowQuiz = async (placeId: string) => {
    if (!session) return;

    const quiz = await getQuiz(session.sessionId, placeId);
    if (quiz) {
      setCurrentQuiz(quiz);
      setShowQuiz(true);
    }
  };

  // クイズ回答
  const handleQuizAnswer = async (
    isCorrect: boolean,
    selectedId: string
  ) => {
    if (!currentQuiz || !session) return;

    // サーバーに送信
    const result = await answerQuiz(session.sessionId, {
      quiz_id: currentQuiz.quiz_id,
      selected_answer_id: selectedId,
    });

    if (result.correct) {
      // 正解：ピース取得
      setPlaceStatuses((prev) => ({
        ...prev,
        [currentQuiz.place_id]: {
          verified: true,
          answered: true,
          correct: true,
        },
      }));

      setShowQuiz(false);
      setCurrentQuiz(null);
      setSelectedPlaceId(null);

      // 全ピース取得チェック
      const allCompleted = FAVORITE_PLACES.every(
        (place) =>
          place.id === currentQuiz.place_id ||
          placeStatuses[place.id]?.correct
      );

      if (allCompleted) {
        // 全ピース取得完了
        setTimeout(async () => {
          await transitionTo('s7');
        }, 1500);
      }
    } else {
      // 不正解：ジャンプスケアと再挑戦
      setShowQuiz(false);
      setShowJumpScare(true);

      setTimeout(() => {
        setShowJumpScare(false);
        // クイズを再表示
        setTimeout(() => {
          setShowQuiz(true);
        }, 500);
      }, 2000);
    }
  };

  // 獲得ピース数
  const obtainedPieces = FAVORITE_PLACES.filter(
    (place) => placeStatuses[place.id]?.correct
  ).length;

  const selectedPlace = FAVORITE_PLACES.find((p) => p.id === selectedPlaceId);

  return (
    <div className={styles.container}>
      {/* タイマー（固定表示） */}
      <Timer
        initialSeconds={timer.seconds}
        large
        fixed
        warningThreshold={TIMER_WARNING_THRESHOLDS.S6}
      />

      {/* コンテンツ */}
      <div className={styles.content}>
        <header className={styles.header}>
          <h1 className={styles.title}>存在証明書の探索</h1>
          <div className={styles.piecesCount}>
            ピース: {obtainedPieces} / {FAVORITE_PLACES.length}
          </div>
        </header>

        {/* 場所リスト */}
        <div className={styles.placesList}>
          {FAVORITE_PLACES.map((place) => {
            const status = placeStatuses[place.id];
            const isObtained = status?.correct;
            const isVerified = status?.verified;
            const isSelected = selectedPlaceId === place.id;

            return (
              <div
                key={place.id}
                className={`${styles.placeItem} ${
                  isObtained ? styles.obtained : ''
                } ${isSelected ? styles.selected : ''}`}
                onClick={() => !isObtained && handleSelectPlace(place.id)}
              >
                <div className={styles.placeIcon}>
                  {isObtained ? '✓' : '📍'}
                </div>
                <div className={styles.placeContent}>
                  <h3 className={styles.placeName}>{place.name}</h3>
                  {isObtained && (
                    <span className={styles.status}>ピース取得済み</span>
                  )}
                  {isVerified && !isObtained && (
                    <span className={styles.status}>到達済み・未回答</span>
                  )}
                </div>
              </div>
            );
          })}
        </div>

        {/* 選択中の場所の操作 */}
        {selectedPlaceId && !placeStatuses[selectedPlaceId]?.verified && (
          <div className={styles.actionPanel}>
            <h3 className={styles.actionTitle}>{selectedPlace?.name}</h3>
            <p className={styles.actionDescription}>
              この場所に到達してください
            </p>

            <div className={styles.actionButtons}>
              <Button
                onClick={handleOpenCamera}
                variant="primary"
                size="large"
                fullWidth
              >
                📷 カメラで撮影
              </Button>

              <Button
                onClick={handleManualVerify}
                variant="secondary"
                size="medium"
                fullWidth
              >
                この場所にいることにする
              </Button>
            </div>
          </div>
        )}

        {/* 到達済み・未回答の場合 */}
        {selectedPlaceId &&
          placeStatuses[selectedPlaceId]?.verified &&
          !placeStatuses[selectedPlaceId]?.correct && (
            <div className={styles.actionPanel}>
              <Button
                onClick={() => loadAndShowQuiz(selectedPlaceId)}
                variant="primary"
                size="large"
                fullWidth
              >
                クイズに挑戦する
              </Button>
            </div>
          )}
      </div>

      {/* カメラモーダル */}
      <Modal
        isOpen={showCamera}
        onClose={() => setShowCamera(false)}
        title="場所を撮影"
        size="large"
      >
        <CameraCapture
          onCapture={handleCaptureComplete}
          onCancel={() => setShowCamera(false)}
        />
      </Modal>

      {/* クイズモーダル */}
      <Modal
        isOpen={showQuiz}
        onClose={() => {}}
        title="あなた自身のことです"
        size="large"
      >
        {currentQuiz && (
          <QuizCard
            question={currentQuiz.question}
            options={currentQuiz.options}
            correctAnswerId={currentQuiz.correct_answer_id}
            onAnswer={handleQuizAnswer}
            placeName={selectedPlace?.name}
          />
        )}
      </Modal>

      {/* ジャンプスケア */}
      {showJumpScare && (
        <JumpScare
          isActive={true}
          text="不正解"
          duration={2000}
          onComplete={() => setShowJumpScare(false)}
        />
      )}
    </div>
  );
}
