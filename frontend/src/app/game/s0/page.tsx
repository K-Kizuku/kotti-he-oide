/**
 * S0: 起動・注意書き・許可取得
 *
 * - カメラ・音声・イヤホンの許可取得
 * - ホラー演出の注意喚起
 * - 言語選択（日本語/English）
 */

'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import styles from './page.module.css';
import Button from '@/components/ui/Button';
import Modal from '@/components/ui/Modal';
import { useSession } from '@/hooks/useSession';

type Language = 'ja' | 'en';

export default function S0Page() {
  const router = useRouter();
  const { session, isLoading: sessionLoading } = useSession();

  const [language, setLanguage] = useState<Language>('ja');
  const [cameraPermission, setCameraPermission] = useState(false);
  const [audioPermission, setAudioPermission] = useState(false);
  const [acceptedWarning, setAcceptedWarning] = useState(false);
  const [showWarningModal, setShowWarningModal] = useState(false);
  const [isRequestingPermissions, setIsRequestingPermissions] = useState(false);

  const texts = {
    ja: {
      title: '赤煉瓦文化館\n〜こっちにおいで〜',
      subtitle: '体験型Webホラーゲーム',
      languageLabel: '言語 / Language',
      requirementsTitle: '必要な準備',
      cameraLabel: 'カメラへのアクセス',
      audioLabel: '音声の再生',
      earphoneLabel: 'イヤホンまたはヘッドホンの装着',
      earphoneNote: '※ 推奨',
      warningTitle: '⚠️ 注意事項',
      warningContent:
        'このゲームにはホラー要素が含まれます。\n驚かせる演出、不快な音声、暗い映像表現があります。\n\n心臓が弱い方、体調の優れない方はプレイをお控えください。',
      acceptWarning: '内容を理解し、同意します',
      requestPermissions: '許可をリクエスト',
      permissionGranted: '✓ 許可済み',
      startGame: 'ゲームを開始',
      permissionError:
        'カメラまたは音声の許可が拒否されました。\nブラウザの設定から許可してください。',
    },
    en: {
      title: 'Red Brick Cultural Hall\n〜Come This Way〜',
      subtitle: 'Interactive Web Horror Game',
      languageLabel: 'Language / 言語',
      requirementsTitle: 'Requirements',
      cameraLabel: 'Camera Access',
      audioLabel: 'Audio Playback',
      earphoneLabel: 'Earphones or Headphones',
      earphoneNote: '※ Recommended',
      warningTitle: '⚠️ Warning',
      warningContent:
        'This game contains horror elements.\nIncludes jump scares, disturbing audio, and dark visuals.\n\nNot recommended for those with heart conditions or feeling unwell.',
      acceptWarning: 'I understand and agree',
      requestPermissions: 'Request Permissions',
      permissionGranted: '✓ Granted',
      startGame: 'Start Game',
      permissionError:
        'Camera or audio permission denied.\nPlease allow access in browser settings.',
    },
  };

  const t = texts[language];

  const requestPermissions = async () => {
    setIsRequestingPermissions(true);

    try {
      // カメラ許可
      const stream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: 'environment' },
        audio: false,
      });

      stream.getTracks().forEach((track) => track.stop());
      setCameraPermission(true);

      // 音声許可（AudioContextで確認）
      const audioContext = new AudioContext();
      await audioContext.resume();
      audioContext.close();
      setAudioPermission(true);
    } catch (error) {
      console.error('Permission error:', error);
      alert(t.permissionError);
    } finally {
      setIsRequestingPermissions(false);
    }
  };

  const canStart =
    cameraPermission && audioPermission && acceptedWarning && session;

  const handleStart = () => {
    if (canStart) {
      router.push('/game/s1');
    }
  };

  if (sessionLoading) {
    return (
      <div className={styles.container}>
        <div className={styles.loading}>Loading...</div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      <div className={styles.content}>
        {/* タイトル */}
        <header className={styles.header}>
          <h1 className={styles.title}>{t.title}</h1>
          <p className={styles.subtitle}>{t.subtitle}</p>
        </header>

        {/* 言語選択 */}
        <section className={styles.section}>
          <label className={styles.label}>{t.languageLabel}</label>
          <div className={styles.languageButtons}>
            <Button
              variant={language === 'ja' ? 'primary' : 'secondary'}
              onClick={() => setLanguage('ja')}
              size="medium"
            >
              日本語
            </Button>
            <Button
              variant={language === 'en' ? 'primary' : 'secondary'}
              onClick={() => setLanguage('en')}
              size="medium"
            >
              English
            </Button>
          </div>
        </section>

        {/* 必要な準備 */}
        <section className={styles.section}>
          <h2 className={styles.sectionTitle}>{t.requirementsTitle}</h2>

          <div className={styles.requirements}>
            <div className={styles.requirement}>
              <span className={styles.requirementIcon}>📷</span>
              <span className={styles.requirementText}>{t.cameraLabel}</span>
              {cameraPermission && (
                <span className={styles.granted}>{t.permissionGranted}</span>
              )}
            </div>

            <div className={styles.requirement}>
              <span className={styles.requirementIcon}>🔊</span>
              <span className={styles.requirementText}>{t.audioLabel}</span>
              {audioPermission && (
                <span className={styles.granted}>{t.permissionGranted}</span>
              )}
            </div>

            <div className={styles.requirement}>
              <span className={styles.requirementIcon}>🎧</span>
              <span className={styles.requirementText}>
                {t.earphoneLabel}
                <span className={styles.note}>{t.earphoneNote}</span>
              </span>
            </div>
          </div>

          {!cameraPermission || !audioPermission ? (
            <Button
              onClick={requestPermissions}
              variant="primary"
              size="large"
              fullWidth
              isLoading={isRequestingPermissions}
            >
              {t.requestPermissions}
            </Button>
          ) : null}
        </section>

        {/* 注意事項 */}
        <section className={styles.section}>
          <div className={styles.warningBox}>
            <h3 className={styles.warningTitle}>{t.warningTitle}</h3>
            <p className={styles.warningText}>{t.warningContent}</p>
          </div>

          <label className={styles.checkbox}>
            <input
              type="checkbox"
              checked={acceptedWarning}
              onChange={(e) => setAcceptedWarning(e.target.checked)}
            />
            <span>{t.acceptWarning}</span>
          </label>
        </section>
      </div>

      {/* 開始ボタン（固定配置） */}
      <div className={styles.footer}>
        <Button
          onClick={handleStart}
          variant="primary"
          size="large"
          fullWidth
          disabled={!canStart}
        >
          {t.startGame}
        </Button>
      </div>

      {/* 警告モーダル */}
      <Modal
        isOpen={showWarningModal}
        onClose={() => setShowWarningModal(false)}
        title={t.warningTitle}
      >
        <p style={{ whiteSpace: 'pre-line' }}>{t.warningContent}</p>
      </Modal>
    </div>
  );
}
