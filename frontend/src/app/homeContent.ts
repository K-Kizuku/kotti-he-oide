export type FeatureCard = {
  id: string;
  icon: string;
  title: string;
  description: string;
  ctaText: string;
  ctaLink: string;
};

export type PWAFeature = {
  id: string;
  name: string;
  description: string;
};

export type HomePageContent = {
  hero: {
    title: string;
    subtitle: string;
    description: string;
  };
  features: FeatureCard[];
  pwaFeatures: PWAFeature[];
  footer: {
    links: { label: string; href: string }[];
    copyright: string;
  };
};

export const homeContent: HomePageContent = {
  hero: {
    title: "Kotti He Oide",
    subtitle: "闇が呼ぶ。こっちに、おいで。",
    description:
      "Web Push通知とカメラフィルターを備えたPWA。深紅のアクセントとノイズに満ちた世界で、確かな体験を。",
  },
  features: [
    {
      id: "push",
      icon: "/icons/notification-96.png",
      title: "Web Push通知",
      description:
        "重要なお知らせをリアルタイムで受け取れます。アプリを閉じていても見逃しません。",
      ctaText: "通知を設定",
      ctaLink: "/notifications",
    },
    {
      id: "camera",
      icon: "/icons/camera-96.png",
      title: "カメラフィルター",
      description:
        "レトロ、ホラー、シリアス、VHS、コミック。5種類のフィルターで世界を変える。",
      ctaText: "フィルターを試す",
      ctaLink: "/camera-filters",
    },
  ],
  pwaFeatures: [
    { id: "offline", name: "オフライン対応", description: "接続が不安定でも動作" },
    { id: "install", name: "ホームに追加", description: "ワンタップで起動" },
    { id: "push", name: "プッシュ通知", description: "重要情報を即時配信" },
    { id: "bg", name: "バックグラウンド同期", description: "見えない所で更新" },
    { id: "app", name: "アプリライク", description: "スムーズな体験" },
    { id: "fast", name: "高速起動", description: "素早く立ち上がる" },
  ],
  footer: {
    links: [
      { label: "通知設定", href: "/notifications" },
      { label: "カメラフィルター", href: "/camera-filters" },
    ],
    copyright: "© 2025 Kotti He Oide App",
  },
};

