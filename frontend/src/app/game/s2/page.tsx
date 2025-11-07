/**
 * S2: お気に入りの場所説明（固定5箇所）
 *
 * 場所:
 * 1. 螺旋階段を見上げる高い天井
 * 2. メインホールの暖炉のレンガ
 * 3. 裏玄関の扉の蝶番
 * 4. 入口エントランスの扉
 * 5. 階上応接室のピアノ
 */

'use client';

import { useState } from 'react';
import styles from './page.module.css';
import Button from '@/components/ui/Button';
import { FAVORITE_PLACES } from '@/lib/game/constants';
import { useGameFlow } from '@/hooks/useGameFlow';

export default function S2Page() {
  const { transitionTo } = useGameFlow();
  const [currentIndex, setCurrentIndex] = useState(0);

  const currentPlace = FAVORITE_PLACES[currentIndex];
  const isLastPlace = currentIndex === FAVORITE_PLACES.length - 1;

  const handleNext = async () => {
    if (isLastPlace) {
      await transitionTo('s3');
    } else {
      setCurrentIndex((prev) => prev + 1);
    }
  };

  const handlePrevious = () => {
    if (currentIndex > 0) {
      setCurrentIndex((prev) => prev - 1);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        {/* ヘッダー */}
        <header className={styles.header}>
          <h1 className={styles.title}>お気に入りの場所</h1>
          <p className={styles.subtitle}>
            この建物には、思い出の詰まった5つの場所があります
          </p>
        </header>

        {/* 進捗インジケーター */}
        <div className={styles.progress}>
          {FAVORITE_PLACES.map((_, index) => (
            <div
              key={index}
              className={`${styles.progressDot} ${
                index === currentIndex ? styles.active : ''
              } ${index < currentIndex ? styles.completed : ''}`}
            />
          ))}
        </div>

        {/* 場所カード */}
        <div className={styles.placeCard}>
          <div className={styles.placeNumber}>
            場所 {currentIndex + 1} / {FAVORITE_PLACES.length}
          </div>

          {/* 場所の画像（未実装の場合はプレースホルダー） */}
          <div className={styles.imageContainer}>
            {currentPlace.imageUrl ? (
              <img
                src={currentPlace.imageUrl}
                alt={currentPlace.name}
                className={styles.placeImage}
              />
            ) : (
              <div className={styles.placeholderImage}>
                <span className={styles.placeholderIcon}>📍</span>
              </div>
            )}
          </div>

          {/* 場所の情報 */}
          <div className={styles.placeInfo}>
            <h2 className={styles.placeName}>{currentPlace.name}</h2>
            <p className={styles.placeDescription}>
              {currentPlace.description}
            </p>
          </div>
        </div>

        {/* ナビゲーションボタン */}
        <div className={styles.navigation}>
          {currentIndex > 0 && (
            <Button
              onClick={handlePrevious}
              variant="secondary"
              size="medium"
            >
              ← 前の場所
            </Button>
          )}

          <div className={styles.spacer} />

          <Button
            onClick={handleNext}
            variant="primary"
            size="medium"
          >
            {isLastPlace ? '次へ進む' : '次の場所 →'}
          </Button>
        </div>
      </div>
    </div>
  );
}
