import BackgroundEffects from "@/components/home/BackgroundEffects";
import HeroSection from "@/components/home/HeroSection";
import FeatureCard from "@/components/home/FeatureCard";
import CTAButton from "@/components/home/CTAButton";
import { homeContent } from "./homeContent";
import styles from "./page.module.css";

export default function Home() {
  const { hero, features, pwaFeatures, footer } = homeContent;

  return (
    <div className={styles.darkRoot}>
      <BackgroundEffects enableNoise enableScanlines enableVignette />

      <main className={styles.main}>
        <HeroSection
          title={hero.title}
          subtitle={hero.subtitle}
          description={hero.description}
        />

        {/* Features */}
        <section className={styles.features} aria-label="主要機能">
          {features.map((f) => (
            <FeatureCard key={f.id} {...f} />
          ))}
        </section>

        {/* PWA 特長 */}
        <section className={styles.pwa} aria-label="PWA機能">
          <h2 className={styles.pwaTitle}>PWAの主な特長</h2>
          <ul className={styles.pwaGrid}>
            {pwaFeatures.map((p) => (
              <li key={p.id} className={styles.pwaItem}>
                <div className={styles.pwaName}>{p.name}</div>
                <div className={styles.pwaDesc}>{p.description}</div>
              </li>
            ))}
          </ul>
          <div className={styles.pwaCtas}>
            <CTAButton text="通知を設定" href="/notifications" />
            <CTAButton text="フィルターを試す" href="/camera-filters" variant="secondary" />
          </div>
        </section>
      </main>

      <footer className={styles.footer}>
        <nav className={styles.footerNav} aria-label="フッターナビゲーション">
          {footer.links.map((l) => (
            <a key={l.href} href={l.href} className={styles.footerLink}>
              {l.label}
            </a>
          ))}
        </nav>
        <p className={styles.copyright}>{footer.copyright}</p>
      </footer>
    </div>
  );
}
