'use client';

/**
 * ErrorModal - ゲームエラー表示専用モーダル
 *
 * GameFlowManagerから発生するエラーを統一的に表示します
 */

import { useEffect } from 'react';
import Modal from '@/components/ui/Modal';
import Button from '@/components/ui/Button';
import { useGameFlow } from '@/hooks/useGameFlow';
import styles from './ErrorModal.module.css';

/**
 * ErrorModal コンポーネント
 *
 * useGameFlowから現在のエラー状態を取得し、自動的にモーダル表示します
 */
export default function ErrorModal() {
  const { currentError, clearError } = useGameFlow();

  // エラーが発生した時のログ出力
  useEffect(() => {
    if (currentError) {
      console.error('[ErrorModal] エラー表示:', currentError);
    }
  }, [currentError]);

  // エラーがない場合は何も表示しない
  if (!currentError) {
    return null;
  }

  // エラーの重大度に応じた処理
  const isRecoverable = currentError.recoverable;

  const handleClose = () => {
    clearError();
  };

  const handleReload = () => {
    if (typeof window !== 'undefined') {
      window.location.reload();
    }
  };

  return (
    <Modal
      isOpen={true}
      onClose={isRecoverable ? handleClose : undefined}
      title="エラー"
      closeOnBackdropClick={isRecoverable}
      showCloseButton={isRecoverable}
      size="medium"
      footer={
        <div className={styles.footer}>
          {isRecoverable ? (
            <Button onClick={handleClose} variant="primary" fullWidth>
              OK
            </Button>
          ) : (
            <>
              <Button onClick={handleReload} variant="primary" fullWidth>
                再読み込み
              </Button>
              <Button
                onClick={() => (window.location.href = '/game/s0')}
                variant="secondary"
                fullWidth
              >
                最初から
              </Button>
            </>
          )}
        </div>
      }
    >
      <div className={styles.errorContent}>
        {/* エラーアイコン */}
        <div className={styles.errorIcon}>⚠️</div>

        {/* ユーザー向けメッセージ */}
        <p className={styles.userMessage}>{currentError.userMessage}</p>

        {/* 開発用情報（本番では非表示にする） */}
        {process.env.NODE_ENV === 'development' && (
          <details className={styles.details}>
            <summary>詳細情報（開発用）</summary>
            <div className={styles.debugInfo}>
              <p>
                <strong>エラーコード:</strong> {currentError.code}
              </p>
              <p>
                <strong>メッセージ:</strong> {currentError.message}
              </p>
              <p>
                <strong>復帰可能:</strong> {isRecoverable ? 'はい' : 'いいえ'}
              </p>
            </div>
          </details>
        )}
      </div>
    </Modal>
  );
}
