import Image from "next/image";
import CTAButton from "./CTAButton";
import styles from "./FeatureCard.module.css";

export type FeatureCardProps = {
  icon: string; // public 配下のパス or emoji
  title: string;
  description: string;
  ctaText: string;
  ctaLink: string;
};

export default function FeatureCard({ icon, title, description, ctaText, ctaLink }: FeatureCardProps) {
  const isEmoji = icon.length <= 3; // お手軽判定
  return (
    <article className={styles.card}>
      <div className={styles.header}>
        {isEmoji ? (
          <div className={styles.emoji} aria-hidden>{icon}</div>
        ) : (
          <div className={styles.iconWrap} aria-hidden>
            <Image src={icon} alt="" width={40} height={40} />
          </div>
        )}
        <h3 className={styles.title}>{title}</h3>
      </div>
      <p className={styles.description}>{description}</p>
      <div className={styles.actions}>
        <CTAButton text={ctaText} href={ctaLink} />
      </div>
    </article>
  );
}

