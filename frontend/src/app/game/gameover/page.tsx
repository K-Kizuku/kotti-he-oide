/**
 * ゲームオーバーページ
 *
 * - タイムアウトなどでゲームが終了した場合に表示
 * - 再挑戦またはトップへ戻る
 */

'use client';

import styles from './page.module.css';
import Button from '@/components/ui/Button';
import GlitchText from '@/components/horror/GlitchText';
import { useGameFlow } from '@/hooks/useGameFlow';

export default function GameOverPage() {
  const { transitionTo } = useGameFlow();

  const handleRetry = async () => {
    await transitionTo('s0');
  };

  const handleGoHome = () => {
    window.location.href = '/';
  };

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        <div className={styles.gameOverText}>
          <GlitchText intensity="high">
            <span className={styles.title}>GAME OVER</span>
          </GlitchText>
        </div>

        <div className={styles.message}>
          <p>時間切れです</p>
          <p>あなたの存在は消えてしまいました...</p>
        </div>

        <div className={styles.actions}>
          <Button onClick={handleRetry} variant="primary" size="large" fullWidth>
            もう一度挑戦する
          </Button>

          <Button
            onClick={handleGoHome}
            variant="secondary"
            size="medium"
            fullWidth
          >
            トップへ戻る
          </Button>
        </div>
      </div>
    </div>
  );
}
