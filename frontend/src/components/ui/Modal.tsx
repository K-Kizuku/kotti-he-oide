/**
 * Modal - モーダルダイアログコンポーネント
 *
 * 警告、確認、情報表示などに使用
 */

'use client';

import { ReactNode, useEffect } from 'react';
import styles from './Modal.module.css';
import Button from './Button';

interface ModalProps {
  /** モーダルを表示するか */
  isOpen: boolean;
  /** モーダルを閉じる関数 */
  onClose?: () => void;
  /** タイトル */
  title?: string;
  /** 子要素 */
  children: ReactNode;
  /** フッター（カスタムボタンなど） */
  footer?: ReactNode;
  /** 閉じるボタンを表示するか */
  showCloseButton?: boolean;
  /** 背景クリックで閉じるか */
  closeOnBackdropClick?: boolean;
  /** サイズ */
  size?: 'small' | 'medium' | 'large';
}

export default function Modal({
  isOpen,
  onClose,
  title,
  children,
  footer,
  showCloseButton = true,
  closeOnBackdropClick = true,
  size = 'medium',
}: ModalProps) {
  // Escキーで閉じる
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && onClose) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
      // スクロール防止
      document.body.style.overflow = 'hidden';
    }

    return () => {
      document.removeEventListener('keydown', handleEscape);
      document.body.style.overflow = '';
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  const handleBackdropClick = () => {
    if (closeOnBackdropClick && onClose) {
      onClose();
    }
  };

  return (
    <div className={styles.backdrop} onClick={handleBackdropClick}>
      <div
        className={`${styles.modal} ${styles[size]}`}
        onClick={(e) => e.stopPropagation()}
        role="dialog"
        aria-modal="true"
        aria-labelledby={title ? 'modal-title' : undefined}
      >
        {title && (
          <div className={styles.header}>
            <h2 id="modal-title" className={styles.title}>
              {title}
            </h2>
            {showCloseButton && onClose && (
              <button
                className={styles.closeButton}
                onClick={onClose}
                aria-label="閉じる"
              >
                ×
              </button>
            )}
          </div>
        )}

        <div className={styles.content}>
          {children}
        </div>

        {footer && (
          <div className={styles.footer}>
            {footer}
          </div>
        )}

        {!footer && showCloseButton && onClose && !title && (
          <div className={styles.footer}>
            <Button onClick={onClose} variant="secondary" fullWidth>
              閉じる
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
