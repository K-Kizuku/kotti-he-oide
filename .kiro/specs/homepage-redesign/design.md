# Design Document

## Overview

「kotti-he-oide」アプリのルートページを、参考LP（赤煉瓦文化館）のダークテーマ・ホラー風デザインを基にした現代的なランディングページに変更する。ダークカラーパレット、グリッチ効果、ノイズテクスチャ、スキャンライン効果を活用して、視覚的にインパクトのあるユーザー体験を提供する。

## Architecture

### デザインシステム
- **カラーパレット**: 参考LPと同様の配色
  - 背景: `#07070a` (ダークブラック)
  - フォアグラウンド: `#e6e6e6` (ライトグレー)
  - アクセント: `#9e0000` (深紅)
  - ミュート: `#0d0d12` (ダークグレー)
- **フォント**: Noto Serif JP（日本語セリフ）をメインに使用
- **視覚効果**: ノイズ、スキャンライン、グリッチ、パルス（心拍）アニメーション

### レイアウト構成
1. **Hero Section**: メインタイトルとキャッチコピー
2. **Features Section**: 主要機能の紹介カード
3. **PWA Features**: PWA機能の説明
4. **Footer**: ナビゲーションとコピーライト

## Components and Interfaces

### 1. Hero Section
```typescript
interface HeroSectionProps {
  title: string;
  subtitle: string;
  description: string;
}
```

**デザイン特徴:**
- 中央配置のタイトル「Kotti He Oide」にグリッチ効果
- サブタイトルにノイズテキスト効果
- 背景にラジアルグラデーション（血飛沫風）
- ビネット効果で画面周囲を暗く

### 2. Feature Cards
```typescript
interface FeatureCardProps {
  icon: string;
  title: string;
  description: string;
  ctaText: string;
  ctaLink: string;
}
```

**デザイン特徴:**
- 半透明の黒背景 `rgba(0, 0, 0, 0.85)`
- 赤色のボーダー `rgba(158, 0, 0, 0.3)`
- バックドロップフィルター（ブラー効果）
- ホバー時のグロー効果強化

### 3. CTA Buttons
```typescript
interface CTAButtonProps {
  text: string;
  href: string;
  variant: 'primary' | 'secondary';
}
```

**デザイン特徴:**
- プライマリ: 赤色背景 `#9e0000`、パルスアニメーション
- セカンダリ: 透明背景、赤色ボーダー
- ホバー時のスケール変化とグロー効果

### 4. Background Effects
```typescript
interface BackgroundEffectsProps {
  enableNoise: boolean;
  enableScanlines: boolean;
  enableVignette: boolean;
}
```

**視覚効果:**
- **ノイズ**: SVGフィルターによるフラクタルノイズ
- **スキャンライン**: CSS repeating-linear-gradientによるCRT効果
- **ビネット**: radial-gradientによる画面周囲の暗化

## Data Models

### Page Content Model
```typescript
interface HomePageContent {
  hero: {
    title: string;
    subtitle: string;
    description: string;
  };
  features: FeatureCard[];
  pwaFeatures: PWAFeature[];
  footer: {
    links: NavigationLink[];
    copyright: string;
  };
}

interface FeatureCard {
  id: string;
  icon: string;
  title: string;
  description: string;
  ctaText: string;
  ctaLink: string;
}

interface PWAFeature {
  id: string;
  name: string;
  description: string;
}
```

## Error Handling

### CSS Fallbacks
- `@supports`クエリでアニメーション対応チェック
- `prefers-reduced-motion`でアニメーション無効化対応
- フォント読み込み失敗時のフォールバック

### Performance Considerations
- CSS animationsの`will-change`プロパティ使用
- 重いエフェクトの条件付き適用
- レスポンシブ対応でモバイル最適化

## Testing Strategy

### Visual Regression Testing
- 各ブレークポイントでのレイアウト確認
- ダークテーマの色彩コントラスト検証
- アニメーション効果の動作確認

### Accessibility Testing
- キーボードナビゲーション対応
- スクリーンリーダー対応
- カラーコントラスト比の検証（WCAG 2.1 AA準拠）

### Browser Compatibility
- モダンブラウザでのCSS Grid/Flexbox対応
- CSS Custom Properties対応
- CSS Animation Timeline対応状況

## Implementation Details

### CSS Architecture
```css
/* カスタムプロパティでテーマ管理 */
:root {
  --background: #07070a;
  --foreground: #e6e6e6;
  --accent: #9e0000;
  --muted: #0d0d12;
  --ring: #6b0000;
}

/* グリッチ効果 */
@keyframes glitchShift {
  0% { transform: translate(0,0); }
  20% { transform: translate(1px,-1px); }
  40% { transform: translate(-1px,1px); }
  60% { transform: translate(2px,0); }
  80% { transform: translate(-2px,0); }
  100% { transform: translate(0,0); }
}

/* パルス効果 */
@keyframes pulse {
  0%, 100% { 
    transform: scale(1); 
    box-shadow: 0 0 0 0 rgba(158,0,0,0.4); 
  }
  50% { 
    transform: scale(1.03); 
    box-shadow: 0 0 0 12px rgba(158,0,0,0); 
  }
}
```

### Component Structure
```
src/app/
├── page.tsx (メインページ)
├── globals.css (グローバルスタイル)
└── components/
    ├── HeroSection.tsx
    ├── FeatureCard.tsx
    ├── CTAButton.tsx
    └── BackgroundEffects.tsx
```

### Responsive Design
- **Desktop**: フル機能、全エフェクト適用
- **Tablet**: 一部エフェクト軽減、レイアウト調整
- **Mobile**: エフェクト最小化、タッチ操作最適化

## Design Decisions and Rationales

### 1. ダークテーマ採用
**決定**: 参考LPと同様の深いダークテーマを採用
**理由**: 
- 現代的で洗練された印象
- 赤色アクセントが映える
- PWAアプリとしての先進性を表現

### 2. グリッチ効果の使用
**決定**: タイトルとキーテキストにグリッチ効果を適用
**理由**:
- 技術的な先進性を視覚的に表現
- ユーザーの注意を引く効果的な演出
- ブランドの独自性を強調

### 3. 日本語フォントの選択
**決定**: Noto Serif JPをメインフォントに採用
**理由**:
- 日本語コンテンツの可読性向上
- セリフ体による上品で落ち着いた印象
- Googleフォントによる安定した配信

### 4. アニメーション戦略
**決定**: パフォーマンスを考慮した軽量アニメーション
**理由**:
- ユーザー体験の向上
- モバイル端末での動作保証
- アクセシビリティ配慮（motion-reduce対応）