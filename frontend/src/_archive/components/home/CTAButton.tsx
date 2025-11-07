import Link from "next/link";
import styles from "./CTAButton.module.css";

export type CTAButtonProps = {
  text: string;
  href: string;
  variant?: "primary" | "secondary";
};

export default function CTAButton({ text, href, variant = "primary" }: CTAButtonProps) {
  return (
    <Link
      href={href}
      className={`${styles.button} ${
        variant === "primary" ? styles.primary : styles.secondary
      }`}
      aria-label={text}
    >
      {text}
    </Link>
  );
}

