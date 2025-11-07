import styles from "./HeroSection.module.css";

export type HeroSectionProps = {
  title: string;
  subtitle: string;
  description: string;
};

export default function HeroSection({ title, subtitle, description }: HeroSectionProps) {
  return (
    <section className={styles.hero} aria-label="アプリ紹介">
      <div className={styles.heroInner}>
        <h1 className={styles.title}>
          <span className={styles.glitch} aria-hidden>
            {title}
          </span>
          <span className={styles.titleText}>{title}</span>
        </h1>
        <p className={styles.subtitle} data-text={subtitle}>
          {subtitle}
        </p>
        <p className={styles.description}>{description}</p>
      </div>
    </section>
  );
}

