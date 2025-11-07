/**
 * ProgressBar - 進捗表示コンポーネント
 *
 * S4の質問進捗などで使用
 */

import styles from './ProgressBar.module.css';

interface ProgressBarProps {
  /** 現在の進捗（0-100） */
  value: number;
  /** 最大値（デフォルト100） */
  max?: number;
  /** ラベル表示 */
  label?: string;
  /** 色のバリエーション */
  variant?: 'default' | 'success' | 'warning' | 'danger';
  /** サイズ */
  size?: 'small' | 'medium' | 'large';
}

export default function ProgressBar({
  value,
  max = 100,
  label,
  variant = 'default',
  size = 'medium',
}: ProgressBarProps) {
  const percentage = Math.min(100, Math.max(0, (value / max) * 100));

  return (
    <div className={styles.container}>
      {label && (
        <div className={styles.label}>
          {label}
        </div>
      )}
      <div className={`${styles.track} ${styles[size]}`}>
        <div
          className={`${styles.fill} ${styles[variant]}`}
          style={{ width: `${percentage}%` }}
          role="progressbar"
          aria-valuenow={value}
          aria-valuemin={0}
          aria-valuemax={max}
        >
          <span className={styles.percentage}>{Math.round(percentage)}%</span>
        </div>
      </div>
    </div>
  );
}
