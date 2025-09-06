# カメラ・フィルターデモ 実装メモ

本メモは `/frontend` の Next.js アプリに実装した「カメラ入力へフィルター適用」デモの引き継ぎ資料です。現状の構成、実装詳細、拡張手順、運用・今後の方針をまとめています。

## 目的 / 概要
- ブラウザのカメラ映像を取得し、Canvas 2D でリアルタイムに画素処理（ポストエフェクト）を行って表示するデモ。
- 5 種類の雰囲気（レトロ/ホラー/シリアス/VHS/コミック）を UI で切り替え可能。
- 元映像（`<video>`）と処理後（`<canvas>`）を並べて表示。

## 動作環境
- Next.js 15（`/frontend`）
- TypeScript、CSS Modules（2 スペース）
- ローカル開発は `http://localhost:3000` で OK。カメラ利用は基本 HTTPS 必須だが、`localhost` は例外で許可されるブラウザが多い。
- iOS Safari ではユーザー操作に紐づいた `getUserMedia` 開始が必要（本実装はボタン押下で開始）。

## セットアップ / 起動
```sh
cd frontend
corepack enable && pnpm i
pnpm dev
# ブラウザで http://localhost:3000/camera-filters
```

## 画面パス / 主要ファイル
- ルート: `/camera-filters`
- ファイル:
  - `src/app/camera-filters/page.tsx` … ページ本体。カメラ起動/停止、描画ループ、UI。
  - `src/app/camera-filters/filters.ts` … フィルター群（Canvas ImageData ベース）。
  - `src/app/camera-filters/noise.ts` … ノイズ群 + ノイズエンジン（複数/詳細設定/持続対応）。
  - `src/app/camera-filters/page.module.css` … デモ用スタイル。

## 実装詳細
- クライアントコンポーネント: `"use client"` を明示。
- カメラ取得: `navigator.mediaDevices.getUserMedia({ video: { facingMode }, audio: false })`
  - フロント/バック切替は `facingMode: "user" | "environment"` を選択。
  - 失敗時はエラーメッセージ表示。停止・再開で復帰可。
- 表示サイズと FPS:
  - 最大横幅 `MAX_WIDTH = 640`、目標 `TARGET_FPS = 24` を `page.tsx` 冒頭で定義。
  - `requestAnimationFrame` で描画し、フレーム間隔でゲーティング（24fps 目安）。
- 描画パイプライン:
  1. `<video>` に MediaStream をアタッチし `play()`。
  2. `<canvas>` に `drawImage(video)`。
  3. `getImageData` → `applyFilter(ImageData, id)` で画素処理。
  4. `putImageData` で反映。
  - 2D コンテキストは `willReadFrequently: true` を指定（頻繁な `getImageData` 読み出し向け）。
- クリーンアップ/安全性:
  - アンマウントや「停止」時に `requestAnimationFrame` を `cancelAnimationFrame`。
  - `MediaStreamTrack.stop()` でカメラ停止。
  - `loadedmetadata` 待ちで `videoWidth/Height` 取得の race を回避。
- UI:
  - 「カメラ開始/停止」ボタン、フィルター選択。
  - ノイズ（複数ノイズを同時適用可能）
    - 各ノイズごとに「有効」「頻度」「規模」「持続f（フレーム数）」「サイズ」を調整。
  - カメラ向き選択（フロント/バック、停止中のみ変更可）。
  - 元映像と処理後を 2 カラムで表示（モバイルは 1 カラム）。

### 実装済みフィルター（filters.ts）
各フィルターは `ImageData` を直接操作します。可能な限り in-place で処理し、近傍参照が必要な場合のみ元フレームコピーを使用。

- `retro`（レトロ）
  - セピア変換、粒子ノイズ（微小ランダム）、ビネット、走査線でフィルム風。
- `horror`（ホラー）
  - 低彩度化 → 緑被り寄せ → 強コントラスト。暗部にランダムノイズで不穏さ。
- `serious`（シリアス）
  - 低彩度 + コントラスト強化 + わずかなガンマ強調 + 弱ビネットでノワール調。
- `vhs`（VHS/グリッチ）
  - R/B チャンネルの水平シフト、走査線、ランダム 1 ラインのグリッチ、わずかな彩度低下。
- `comic`（コミック）
  - ポスタライズ（段階化） + Sobel で輪郭抽出し、閾値以上を黒で上描き。

補助処理:
- `applyVignette`（周辺減光）、`addScanlines`（走査線）、`adjustContrast`（コントラスト調整）、`toGray`（輝度計算）など。

### 実装済みノイズ（noise.ts）
フィルター処理後に「別レイヤーの偶発的破損」として適用されます。複数ノイズを同時に使用可能。ノイズエンジンが各タイプに対して「イベント（位置/サイズ/モード/残存フレーム）」を生成・管理し、`durationFrames` の間だけ同一領域を乱します。

ノイズ共通パラメータ（`NoiseParams`）
- `enabled` … 有効/無効
- `frequency` … 0..1（フレーム毎の発生確率。例: 0.2 で約20%）
- `magnitude` … 0..1（影響の強さ/振幅）
- `durationFrames` … 1..60（持続フレーム数）
- `size` … 0..1（帯幅/ブロック/パッチ等の幾何サイズ）

推奨値の目安
- `frequency`: 0.05〜0.3（高すぎると常時ノイズ）
- `durationFrames`: 1〜8（ブラーや残像感不要なら 1〜3）
- `size`: 0.1〜0.4（画面の 10〜40% 程度の帯/ブロック/パッチ）

- `dropout`（ドロップアウト）
  - 細い水平ラインの白飛び/黒潰れ/近傍コピーによる欠損をランダム発生。
- `block`（ブロック欠損/ずれ）
  - 矩形ブロックを別位置からコピーして欠損/ずれを表現。
- `tear`（ティア）
  - 水平帯域が横方向に波打つ/ずれる時間軸エラー風の乱れ。
- `snow`（スノーノイズ）
  - 局所矩形パッチを白黒ランダム値で砂嵐化。
- `headswitch`（ヘッドスイッチング）
  - 下端にカラーストライプと乱れを周期的に発生（VHS 的）。

適用順序:
1. フィルター `applyFilter(frame, filterId)`
2. ノイズエンジン `noiseEngine.step(frame)`（複数ノイズ・持続・確率・サイズを総合適用）

## 拡張手順（新規フィルターの追加）
1. `filters.ts` の `FilterId` に新しい ID を追加。
2. `FILTERS` 配列へ `{ id, label, description }` を追加（UI へ自動反映）。
3. `applyYourFilter(image: ImageData)` を実装。
4. `applyFilter` の `switch` にケースを追加。
5. 必要に応じて `page.tsx` に UI（スライダー等）を追加し、関数引数に渡す。

テンプレート例:
```ts
// filters.ts
export type FilterId = "retro" | "horror" | "serious" | "vhs" | "comic" | "your";
export function applyYour(image: ImageData) {
  const { data, width: w, height: h } = image;
  for (let i = 0; i < data.length; i += 4) {
    const r = data[i], g = data[i+1], b = data[i+2];
    // ここに画素処理
    data[i] = r; data[i+1] = g; data[i+2] = b;
  }
}
```

## 拡張手順（新規ノイズの追加）
1. `noise.ts` の `NoiseId` に ID を追加。
2. `NOISES` 配列へ `{ id, label, description }` を追加（UI へ自動反映）。
3. エンジンにイベント型を追加し、`spawn()`（初期位置/サイズ/モード生成）と `applyEvents()`（適用/TTL減算）を拡張。
4. `defaultNoiseConfig()` にデフォルトパラメータを追加。
5. UI（`page.tsx`）は `NOISES` に基づき自動で行が増えるため、通常変更不要（追加パラメータが必要な場合のみUIを拡張）。

テンプレート例:
```ts
export type NoiseId = "dropout" | "block" | "tear" | "snow" | "headswitch" | "your";
function noiseYour(image: ImageData, strength: number) {
  const { data, width: w, height: h } = image;
  // strength は 0..100、頻度/規模に利用
  // 例: if (Math.random() > strength/100) return; // 発生確率
  // 部分領域だけを処理するのがポイント
}
```

## テスト/検証のヒント
- 体感でカクつきが出る場合は `MAX_WIDTH` を 480 へ下げる、`TARGET_FPS` を 20 に調整。
- モバイルでフロントカメラがミラー表示されない点が気になる場合は、CSS の `transform: scaleX(-1)` を適用（元映像だけ/出力だけ等、要件に応じて）。
- E2E（任意）: Playwright でページ遷移・UI 操作まで（カメラ許可は CI だと制約あり）。

## 既知の注意点
- 低スペック端末では `comic`（Sobel）や `vhs`（色ずれ + コピー）で負荷が上がりやすい。
- `Math.random()` によるノイズはフリッカーが出る。必要であればシード付き PRNG に差し替え可能。
- `getUserMedia` はユーザー操作トリガー必須の環境あり。サイレント起動は失敗する場合がある。
- ブラウザ差異（iOS Safari、旧 Android WebView など）で `facingMode: "environment"` が無視されることがある（その際はデフォルトのカメラが使われる）。

## コード規約の準拠
- TypeScript / Next.js App Router / CSS Modules。
- 2 スペースインデント、ESLint（`pnpm lint`）。
- 公開値は `NEXT_PUBLIC_` のみ使用（本機能では不要）。

## 今後の実装方針（推奨ロードマップ）
- パフォーマンス最適化
  - `requestVideoFrameCallback` の利用（対応ブラウザでは映像フレーム同期/ドロップ制御が容易）。
  - Web Worker + `OffscreenCanvas` による処理オフロード（メインスレッドの UI 応答性確保）。
  - WebGL/WebGPU 版フィルター（Fragment Shader での高速化、より高度な効果）。
  - メモリアロケーションの削減（`ImageData` の再利用、バッファプーリング）。
- 機能拡張
  - フィルター強度や各種パラメータを UI スライダーで調整（彩度/コントラスト/ビネット量/閾値など）。
  - ノイズの多重適用（済）と、発生頻度・継続時間・領域サイズの個別パラメータ（済）。
  - 複数フィルターの同時比較グリッド表示。
  - 静止画スナップショット保存、動画録画（`MediaRecorder`）とダウンロード。
  - 自撮り用のミラー ON/OFF トグル、アスペクト比/解像度の選択肢追加。
- 設計整備
  - フィルターをプラグイン化（インターフェイス/登録メカニズム整備、ドキュメント化）。
  - 型安全なパラメータ定義（各フィルターの `schema` と UI 自動生成）。
  - ベンチマーク用ユーティリティ（平均処理時間、ドロップ率の計測）。
- QA/運用
  - モバイル主要機種での実機確認（iOS Safari/Chrome Android）。
  - 許可ダイアログ/ブロック時の導線強化（モーダルで手順案内）。

## よくあるトラブルと対処
- カメラが起動しない
  - ブラウザ設定でカメラ許可を確認。`https` 以外は拒否される環境もある（`localhost` は多くのブラウザで許可）。
- 画面が真っ暗/サイズが 0
  - `loadedmetadata` 前にサイズ参照している可能性 → 実装済の待機ロジックを維持すること。
- 映像が重い
  - `MAX_WIDTH`/`TARGET_FPS` を下げる。Sobel 等の高コストフィルターを避ける/間引く。

## 保守の観点
- 依存が薄い（DOM API + Canvas 2D のみ）。ライブラリ更新の影響は限定的。
- 将来 WebGL/WebGPU に移行する場合も現行 UI と切り替え可能な設計にしておくと安全。

---
更新日: 2025-09-06
担当: （記入）
