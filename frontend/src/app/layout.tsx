import type { Metadata } from "next";
import { Geist, Geist_Mono, Noto_Serif_JP } from "next/font/google";
import localFont from "next/font/local";
import "./globals.css";
import { GameFlowProvider } from "@/contexts/GameFlowContext";
import ErrorModal from "@/components/game/ErrorModal";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const notoSerifJP = Noto_Serif_JP({
  variable: "--font-noto-serif-jp",
  subsets: ["latin"],
  weight: ["400", "700"],
});

// 怨霊フォント（ホラー演出用）
const onryouFont = localFont({
  src: "../../public/fonts/onryou.ttf",
  variable: "--font-onryou",
  display: "swap",
});

export const metadata: Metadata = {
  title: "赤煉瓦文化館 〜こっちにおいで〜",
  description: "福岡市赤煉瓦文化館（現エンジニアカフェ）を舞台とした、体験型Webホラーゲーム",
  manifest: "/manifest.json",
  themeColor: "#000000",
  appleWebApp: {
    capable: true,
    statusBarStyle: "black-translucent",
    title: "こっちにおいで",
  },
  icons: {
    icon: "/icons/icon-192.png",
    apple: "/icons/icon-192.png",
  },
  other: {
    "mobile-web-app-capable": "yes",
    "apple-mobile-web-app-capable": "yes",
    "apple-mobile-web-app-status-bar-style": "default",
    "apple-mobile-web-app-title": "KottiApp",
    "msapplication-TileColor": "#000000",
    "msapplication-TileImage": "/icons/icon-144.png",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <head>
        <link rel="manifest" href="/manifest.json" />
        <meta name="theme-color" content="#000000" />
        <link rel="icon" href="/icons/icon-192.png" />
        <link rel="apple-touch-icon" href="/icons/icon-192.png" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="default" />
        <meta name="apple-mobile-web-app-title" content="KottiApp" />
        <meta name="mobile-web-app-capable" content="yes" />
        <meta name="msapplication-TileColor" content="#000000" />
        <meta name="msapplication-TileImage" content="/icons/icon-144.png" />
      </head>
      <body className={`${geistSans.variable} ${geistMono.variable} ${notoSerifJP.variable} ${onryouFont.variable}`}>
        <GameFlowProvider>
          {children}
          <ErrorModal />
        </GameFlowProvider>
      </body>
    </html>
  );
}
