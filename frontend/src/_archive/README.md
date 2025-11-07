# アーカイブディレクトリ

このディレクトリには、技術検証用に作成されたページ・コンポーネント・ライブラリが格納されています。

## 目的

- Next.js App Routerのルーティング対象から除外（`_`プレフィックス）
- 本番で使用するロジックを含む検証コードを保存
- 将来的に参照・再利用可能な状態で保持

**重要:** このディレクトリ内のファイルは削除・変更しないでください。本番で使用するロジックが含まれています。

---

## ディレクトリ構造

```
_archive/
├── app-pages/              # 検証用App Routerページ
│   ├── camera-filters/     # カメラフィルターデモページ
│   └── notifications/      # Web Push通知テストページ
├── components/             # 検証用コンポーネント
│   ├── home/               # ホームページ装飾用コンポーネント
│   └── push/               # Web Push関連UIコンポーネント
├── hooks/                  # Web Push関連カスタムフック
│   ├── useNotificationPermission.ts
│   └── usePushNotification.ts
├── lib/                    # ライブラリ・ユーティリティ
│   ├── api/
│   │   └── pushApi.ts      # Web Push API クライアント
│   └── utils/
│       ├── notificationUtils.ts
│       └── serviceWorker.ts
└── types/                  # 型定義
    └── push.ts             # Web Push型定義
```

---

## 本番で使用する重要なロジック

### 1. カメラフィルターロジック

**場所:** `_archive/app-pages/camera-filters/`

**重要ファイル:**
- `filters.ts` - 5種類のフィルター実装（retro, horror, serious, VHS, comic）
- `noise.ts` - ノイズ生成ユーティリティ
- `page.tsx` - カメラ処理の実装例

**本番での使用例:**
```typescript
import { applyFilter } from '@/_archive/app-pages/camera-filters/filters';
import { generateNoise } from '@/_archive/app-pages/camera-filters/noise';
```

**主な機能:**
- Canvas 2D APIを使った高速ピクセル処理
- リアルタイム映像へのフィルター適用
- ホラー演出用の色味・ノイズ効果

---

### 2. Web Push通知システム

**場所:** `_archive/components/push/`, `_archive/hooks/`, `_archive/lib/`

**重要ファイル:**
- `hooks/usePushNotification.ts` - プッシュ通知フック
- `lib/api/pushApi.ts` - プッシュ通知APIクライアント
- `lib/utils/serviceWorker.ts` - Service Worker管理

**本番での使用例:**
```typescript
import { usePushNotification } from '@/_archive/hooks/usePushNotification';
import { PushAPI } from '@/_archive/lib/api/pushApi';
```

**主な機能:**
- RFC 8030/8291/8292準拠のWeb Push実装
- VAPID認証とメッセージ暗号化
- サブスクリプション管理

---

### 3. ホームページ装飾コンポーネント

**場所:** `_archive/components/home/`

**ファイル:**
- `BackgroundEffects.tsx` - 背景エフェクト
- `CTAButton.tsx` - CTA（Call To Action）ボタン
- `FeatureCard.tsx` - 機能紹介カード
- `HeroSection.tsx` - ヒーローセクション

**本番での使用例:**
```typescript
import { HeroSection } from '@/_archive/components/home/HeroSection';
```

---

## インポートパスについて

このディレクトリ内のファイルは、TypeScriptパスエイリアス（`@/*`）を使ってインポートできます：

```typescript
// ✅ 推奨
import { applyFilter } from '@/_archive/app-pages/camera-filters/filters';

// ❌ 非推奨（相対パス）
import { applyFilter } from '../../../_archive/app-pages/camera-filters/filters';
```

---

## 注意事項

1. **ルーティング対象外:** `_`プレフィックスにより、Next.jsのApp Routerはこのディレクトリをルーティング対象として認識しません
2. **型チェック:** ビルド時に型チェックは実行されるため、エラーがあるとビルドが失敗します
3. **git履歴:** `git mv`で移動しているため、変更履歴は保持されています
4. **依存関係:** このディレクトリ内のファイルは本番コードから参照される可能性があります

---

## 参考リンク

- [Next.js App Router - Private Folders](https://nextjs.org/docs/app/building-your-application/routing/colocation#private-folders)
- [Canvas 2D API - MDN](https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API)
- [Web Push API - MDN](https://developer.mozilla.org/en-US/docs/Web/API/Push_API)
