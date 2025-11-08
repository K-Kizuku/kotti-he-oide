/**
 * ホームページ - QRコードからのランディング
 *
 * ゲーム開始前の案内ページ
 */

import Link from "next/link";
import styles from "./page.module.css";

export default function Home() {
  return (
    <div className={styles.darkRoot}>
      <main className={styles.main}>
        <section style={{ textAlign: "center", padding: "2rem" }}>
          <h1 style={{ fontSize: "2rem", marginBottom: "1rem" }}>
            赤煉瓦文化館 <br />
            〜こっちにおいで〜
          </h1>
          <p style={{ marginBottom: "2rem", color: "#666" }}>
            体験型Webホラーゲーム
          </p>

          <Link
            href="/game/s0"
            style={{
              display: "inline-block",
              padding: "1rem 2rem",
              background: "#000",
              color: "#fff",
              textDecoration: "none",
              borderRadius: "4px",
              fontSize: "1.2rem",
            }}
          >
            ゲームを開始
          </Link>

          <div style={{ marginTop: "3rem", fontSize: "0.9rem", color: "#999" }}>
            <p>プレイ時間：20〜30分（館内移動込み）</p>
            <p>推奨環境：カメラ・音声・イヤホン</p>
          </div>
        </section>
      </main>
    </div>
  );
}
