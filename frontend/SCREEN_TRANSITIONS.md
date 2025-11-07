# 画面遷移仕様書

赤煉瓦文化館ホラーゲーム「こっちにおいで」の全画面説明と遷移条件をまとめたドキュメントです。

## 目次

- [画面遷移フロー図](#画面遷移フロー図)
- [各画面の詳細](#各画面の詳細)
- [特殊な遷移条件](#特殊な遷移条件)
- [エラー処理](#エラー処理)

---

## 画面遷移フロー図

```
START (/)
  ↓
[S0] 起動・注意書き・許可取得
  ↓ 条件: カメラ・音声許可 + 注意事項同意
[S1] 1942年パート - 担当者との会話
  ↓ 条件: 会話完了
[S2] お気に入りの場所説明（5箇所）
  ↓ 条件: 全5箇所確認
[S3] 移動指示 + ホラー演出
  ↓ 条件: 2階診査室への移動指示確認
[S4] 診査室 - 内省10問
  ↓ 条件: 全10問回答完了
[S5] 死亡届受理 + 7分制限開始
  ↓ 条件: S6探索開始ボタン押下
[S6] 存在証明書探索（7分タイマー）
  ↓ 条件: 全5ピース取得 OR タイムアウト
  ├─ 成功 → [S7]
  └─ 失敗 → [GAMEOVER]
[S7] 2002年パート - メッセージ受け取り
  ↓ 条件: S4の「人生の最期に達成したいこと」を完全一致で再入力
[S8] メインホール3分探索
  ↓ 条件: 痕跡発見 OR タイムアウト
  ├─ 成功 → [S9]
  └─ 失敗 → [GAMEOVER]
[S9] 2025年 - メッセージ刻み
  ↓ 条件: メッセージ保存完了
END (/) トップへ戻る

[GAMEOVER]
  ↓ 選択
  ├─ もう一度挑戦 → [S0]
  └─ トップへ戻る → (/)
```

---

## 各画面の詳細

### S0: 起動・注意書き・許可取得

**パス**: `/game/s0`

**目的**:
- ゲーム開始前の準備とプレイヤーへの注意喚起
- カメラ・音声の許可取得
- ホラー要素への同意確認

**画面要素**:
- 言語選択（日本語 / English）
- カメラアクセス許可ボタン
- 音声再生許可
- イヤホン装着推奨表示
- ホラー演出の注意事項チェックボックス
- ゲーム開始ボタン

**遷移条件**:
```typescript
// 以下の全てが true の場合のみ遷移可能
✓ カメラ許可済み (cameraPermission === true)
✓ 音声許可済み (audioPermission === true)
✓ 注意事項に同意 (acceptedWarning === true)
✓ セッションが存在 (session !== null)
```

**遷移先**: `S1`

**実装ファイル**: `src/app/game/s0/page.tsx`

---

### S1: 1942年パート - 担当者との会話

**パス**: `/game/s1`

**目的**:
- プレイヤーに1942年の時代設定を体験させる
- VOICEVOX（青山龍星）による音声で臨場感を演出
- プレイヤー情報の入力

**画面要素**:
- カメラプレビュー（暗いオーバーレイ付き）
- 担当者の音声再生（VOICEVOX）
- 字幕表示
- テキスト入力フォーム
  - 来館方法
  - 普段の活動

**会話ステップ**:
1. `welcome`: ようこそメッセージ
2. `ask_visit_method`: 来館方法を質問
3. `input_visit_method`: プレイヤー入力
4. `ask_activities`: 普段の活動を質問
5. `input_activities`: プレイヤー入力
6. `closing`: 締めの挨拶
7. `complete`: 完了

**遷移条件**:
```typescript
// 全ての会話ステップが完了
step === 'complete'
```

**遷移先**: `S2`

**実装ファイル**: `src/app/game/s1/page.tsx`

---

### S2: お気に入りの場所説明（5箇所）

**パス**: `/game/s2`

**目的**:
- 館内の重要な5箇所をプレイヤーに紹介
- 後のS6探索パートで必要となる場所の事前確認

**表示される場所**:
1. 螺旋階段を見上げる高い天井 (`spiral_stairs`)
2. メインホールの暖炉のレンガ (`fireplace`)
3. 裏玄関の扉の蝶番 (`back_door_hinge`)
4. 入口エントランスの扉 (`entrance_door`)
5. 階上応接室のピアノ (`piano_room`)

**画面要素**:
- 進捗インジケーター（5箇所中何箇所目）
- 場所の画像
- 場所の名前と説明
- 前へ/次へボタン

**遷移条件**:
```typescript
// 最後の場所（5箇所目）で「次へ」ボタンを押下
currentIndex === 4 && isLastPlace === true
```

**遷移先**: `S3`

**実装ファイル**: `src/app/game/s2/page.tsx`

---

### S3: 移動指示 + ホラー演出

**パス**: `/game/s3`

**目的**:
- 2階診査室への移動を促す
- ホラー演出で雰囲気を盛り上げる
- プレイヤーに恐怖感を与える

**画面要素**:
- カメラプレビュー
- ホラーオーバーレイ（ノイズ + ビネット効果）
- ランダムで表示される影・シルエット
- グリッチテキストで指示表示
- 環境音（horror_ambient.mp3）
- 「こっちに来てはならない」音声

**演出タイミング**:
- 3秒後: 影を表示
- 5秒後: 「2階の診査室へお越しください」指示表示
- 8秒ごと: 影の点滅

**遷移条件**:
```typescript
// 「診査室へ進む」ボタンを押下
// 条件なし（いつでも遷移可能）
```

**遷移先**: `S4`

**実装ファイル**: `src/app/game/s3/page.tsx`

---

### S4: 診査室 - 内省10問

**パス**: `/game/s4`

**目的**:
- プレイヤーの過去や価値観を引き出す
- 回答はS6のクイズ生成に使用
- 「人生の最期に達成したいこと」はS7で使用

**固定10問**:
1. 小学生時代に何に夢中でしたか？
2. 子供の頃、尊敬していた人は誰ですか？
3. 初めてやりがいを感じた仕事や活動は何ですか？
4. 人生で一番悔しかったことは何ですか？
5. あなたにとって「成功」とは何ですか？
6. 最も影響を受けた本や映画は何ですか？
7. 人生で最も大切にしている価値観は何ですか？
8. あなたの理想の生き方はどんなものですか？
9. **人生の最期に達成したいことは何ですか？**（S7で使用）
10. あなたの名前を教えてください

**画面要素**:
- プログレスバー（10問中何問目）
- 質問テキスト
- テキストエリア（回答入力）
- 次へボタン

**バリデーション**:
```typescript
// 無効な回答をブロック
✗ 空文字列
✗ "なし"
✗ "特にない"
✗ "無し"
✗ 空白のみ
```

**保存方式**:
- 各回答は逐次サーバーに保存
- Wi-Fi切断やリロード後も復元可能
- LocalStorage + サーバー側の両方で保持

**遷移条件**:
```typescript
// 全10問に有効な回答を入力済み
answeredQuestions.length === 10
```

**遷移先**: `S5`

**実装ファイル**: `src/app/game/s4/page.tsx`

---

### S5: 死亡届受理 + 7分制限開始

**パス**: `/game/s5`

**目的**:
- 1972年パートへの時代遷移を演出
- S6の7分タイマー開始を告知
- プレイヤーに緊張感を与える

**画面要素**:
- 年代表示：「1972年」
- グリッチテキストでメッセージ表示
- 段階的なステップ表示
  1. `notification`: 「死亡届が受理されました」
  2. `explanation`: 「7分以内に存在証明書を見つけてください」
  3. `ready`: 「探索を開始」ボタン表示

**演出タイミング**:
- 3秒後: notification → explanation
- 6秒後: explanation → ready

**遷移条件**:
```typescript
// 「探索を開始」ボタンを押下
// + サーバー側でS6開始時刻を記録
step === 'ready' && ボタン押下
```

**遷移先**: `S6`

**API呼び出し**:
```typescript
// S6開始時刻を記録
await startS6Exploration(sessionId);
```

**実装ファイル**: `src/app/game/s5/page.tsx`

---

### S6: 存在証明書探索（7分タイマー）★ 最重要

**パス**: `/game/s6`

**目的**:
- ゲームの中核となる探索パート
- 5箇所すべてでピースを取得
- 制限時間内にクリアする緊張感

**制限時間**: 7分（420秒）

**タスク**:
各場所で以下を実行:
1. **場所に到達**: カメラ撮影 or Web選択
2. **クイズに回答**: S4の回答を元にした4択クイズ
3. **ピース取得**: 正解するとピース獲得

**5つの場所**:
- `spiral_stairs`: 螺旋階段
- `fireplace`: 暖炉
- `back_door_hinge`: 裏玄関の蝶番
- `entrance_door`: 入口エントランス
- `piano_room`: ピアノ

**到達判定方法**:

**方法1: カメラ撮影**
```typescript
// 画像類似度判定（サーバー側）
similarity >= 0.5 → 到達成功 (verified_by: "photo")
similarity < 0.5 → 到達失敗（再試行可能）
```

**方法2: Web選択**
```typescript
// 場所リストから選択
verified_by: "manual"
```

**クイズシステム**:
- S6入室時に5問すべて事前生成
- 4択形式
- 正解: プレイヤー自身のS4回答
- ダミー選択肢:
  1. プレイヤーの別の回答
  2. 過去プレイヤーの匿名回答
  3. システム汎用回答

**不正解時の処理**:
```typescript
// ジャンプスケア演出
showJumpScare = true (3秒間)
↓
// 即座に再挑戦可能
クイズモーダル再表示
```

**画面要素**:
- タイマー（カウントダウン表示）
- 場所リスト（5箇所、進捗表示）
- カメラキャプチャモーダル
- クイズモーダル
- ジャンプスケア画面

**遷移条件**:

**成功パターン**:
```typescript
// 全5箇所でピース取得完了
allPlacesCompleted = FAVORITE_PLACES.every(
  place => placeStatuses[place.id]?.correct === true
)
→ S7へ遷移
```

**失敗パターン**:
```typescript
// 7分経過時点で未完了
timer.onComplete() && !allPlacesCompleted
→ GAMEOVERへ遷移
```

**特殊ルール**:
- 最後のピース回答中にタイムアウトした場合のみ猶予あり
- 回答送信まで待機してから判定

**遷移先**:
- 成功: `S7`
- 失敗: `GAMEOVER`

**実装ファイル**: `src/app/game/s6/page.tsx`

---

### S7: 2002年パート - メッセージ受け取り

**パス**: `/game/s7`

**目的**:
- 2002年への時代遷移
- S4で答えた重要な質問の確認
- プレイヤー自身の記憶の再確認

**画面要素**:
- 年代表示：「2002年」
- 導入メッセージ
- テキスト入力フォーム
- エラーメッセージ表示

**ステップ**:
1. `intro`: 導入メッセージ表示
2. `input`: 「人生の最期に達成したいこと」を再入力
3. `success`: 成功メッセージ表示

**検証ロジック**:
```typescript
const trimmedInput = userInput.trim();
const trimmedAnswer = lifeGoalAnswer.trim();

if (trimmedInput === trimmedAnswer) {
  // 完全一致 → 成功
  setStep('success');
  setTimeout(() => transitionTo('s8'), 3000);
} else {
  // 不一致 → エラー
  setError('あなたの答えと一致しません。\n一字一句、正確に思い出してください。');
}
```

**重要**:
- **完全一致が必須**（スペース、句読点も含む）
- トリミング後の比較
- 不一致の場合は何度でも再入力可能

**遷移条件**:
```typescript
// S4で回答した「人生の最期に達成したいこと」と完全一致
userInput.trim() === savedAnswer.trim()
→ 3秒後にS8へ自動遷移
```

**遷移先**: `S8`

**実装ファイル**: `src/app/game/s7/page.tsx`

---

### S8: メインホール3分探索

**パス**: `/game/s8`

**目的**:
- メインホールで過去プレイヤーの痕跡を探す
- 制限時間内に発見する

**制限時間**: 3分（180秒）

**画面要素**:
- タイマー（カウントダウン表示）
- マップ画像（該当箇所を探す）
- 「痕跡を発見」ボタン
- 過去プレイヤーのメッセージ一覧

**ステップ**:
1. `search`: 痕跡探索中（3分タイマー）
2. `messages`: 発見成功、メッセージ一覧表示

**遷移条件**:

**成功パターン**:
```typescript
// 「痕跡を発見」ボタンを押下
handleFound() → タイマー停止 → 過去メッセージ取得 → S9へ
```

**失敗パターン**:
```typescript
// 3分経過
timer.onComplete() && step !== 'messages'
→ GAMEOVERへ遷移
```

**表示されるメッセージ**:
- 過去プレイヤーが刻んだメッセージ一覧
- 完全匿名（session_idのみ保持）
- 作成日時順に表示

**遷移先**:
- 成功: `S9`
- 失敗: `GAMEOVER`

**実装ファイル**: `src/app/game/s8/page.tsx`

---

### S9: 2025年 - メッセージ刻み

**パス**: `/game/s9`

**目的**:
- 現在（2025年）への帰還
- プレイヤー自身のメッセージを刻む
- 次のプレイヤーへの継承

**画面要素**:
- 年代表示：「2025年 - 現在」
- メッセージ入力フォーム（最大120文字）
- 文字数カウンター
- 場所選択（5箇所から1箇所）
- 保存完了メッセージ
- 「終了」ボタン

**ステップ**:
1. `intro`: 導入メッセージ
2. `input`: メッセージ入力
3. `select`: 刻む場所を選択
4. `complete`: 保存完了

**メッセージ制約**:
```typescript
// 最大120文字
message.length <= MAX_MESSAGE_LENGTH (120)

// 空でないこと
message.trim().length > 0
```

**保存内容**:
```typescript
{
  session_id: string,      // 匿名化（ハッシュ化推奨）
  message_text: string,    // 最大120文字
  place_id: PlaceId,       // 5箇所のいずれか
  created_at: timestamp    // UTC
}
```

**遷移条件**:
```typescript
// メッセージ保存完了
await saveMessage(sessionId, { message_text, place_id });
→ step = 'complete'
```

**終了方法**:
```typescript
// 「終了」ボタンでトップページへ
window.location.href = '/'
```

**遷移先**: `/` (トップページ)

**実装ファイル**: `src/app/game/s9/page.tsx`

---

### GAMEOVER: ゲームオーバー

**パス**: `/game/gameover`

**目的**:
- タイムアウト時の処理
- 再挑戦またはトップへ戻る選択肢

**到達経路**:
- S6で7分タイムアウト
- S8で3分タイムアウト

**画面要素**:
- グリッチエフェクト付き「GAME OVER」
- 「時間切れです」メッセージ
- 「あなたの存在は消えてしまいました...」
- 2つのボタン:
  - 「もう一度挑戦する」
  - 「トップへ戻る」

**遷移先**:
- もう一度挑戦: `S0`（新規セッションで開始）
- トップへ戻る: `/`

**実装ファイル**: `src/app/game/gameover/page.tsx`

---

## 特殊な遷移条件

### バリデーション付き遷移

#### S6 → S7
```typescript
// GameFlowManager内で自動検証
validateS6ToS7Transition(sessionId) {
  const progress = await getS6Progress(sessionId);

  // 5箇所全てチェック
  if (progress.length !== 5) return false;

  // 全てのピースが取得済みか確認
  const allCompleted = progress.every(p => p.verified && p.correct);

  return allCompleted;
}
```

**エラー時の動作**:
- モーダルでエラーメッセージ表示
- 現在のシーンに留まる
- ユーザーメッセージ: 「まだ全てのピースが揃っていません（X/5）」

#### S7 → S8
```typescript
// S7ページ内で検証
validateS7ToS8Transition(sessionId) {
  const answers = await getAnswers(sessionId);
  const lifeGoalAnswer = answers.find(
    a => a.question_id === '9' || a.question_id === 'q9_life_goal'
  );

  return lifeGoalAnswer !== null;
}
```

**エラー時の動作**:
- エラーメッセージ表示
- 入力フォームで再入力を促す

---

### タイマー連動遷移

#### S6: 7分タイマー
```typescript
const timer = useTimer({
  initialSeconds: 420, // 7分
  autoStart: true,
  warningThreshold: 60, // 残り1分で警告
  onComplete: async () => {
    // 最後のピース回答中でない場合
    if (!isAnsweringLastPiece) {
      await transitionTo('gameover');
    }
  }
});
```

**警告表示**:
- 残り1分: タイマーの色が変わる
- 残り30秒: 点滅開始

#### S8: 3分タイマー
```typescript
const timer = useTimer({
  initialSeconds: 180, // 3分
  autoStart: true,
  warningThreshold: 30, // 残り30秒で警告
  onComplete: async () => {
    if (step !== 'messages') {
      await transitionTo('gameover');
    }
  }
});
```

**猶予なし**:
- S8は猶予がない（S6と異なる）
- タイムアウト時は即座にGAMEOVER

---

### ブラウザバック制御

全ゲームシーン（S0〜S9、GAMEOVER）でブラウザの戻るボタンが無効化されます。

**実装方法**:
```typescript
// GameFlowManager初期化時
browserHistoryBlocker.enable();

// 動作
window.addEventListener('popstate', (e) => {
  e.preventDefault();
  window.history.pushState(null, '', window.location.href);
});
```

**目的**:
- ゲーム進行の保護
- 意図しない遷移の防止
- シーン順序の強制

---

## エラー処理

### エラータイプと対応

#### 1. SESSION_EXPIRED (セッション期限切れ)

**発生条件**:
```typescript
// 60分経過後
new Date() > new Date(session.expiresAt)
```

**動作**:
- ErrorModalで通知
- LocalStorageのセッションをクリア
- 自動的にS0へリダイレクト

**ユーザーメッセージ**:
> セッションの有効期限が切れました。最初からやり直してください。

**復帰可能**: ❌ 不可（新規セッション開始が必要）

---

#### 2. VALIDATION_ERROR (遷移条件未達)

**発生例**:
- S6で全ピース未取得なのにS7へ遷移しようとした
- S7で回答が不一致

**動作**:
- ErrorModalで通知
- 現在のシーンに留まる
- 再試行可能

**ユーザーメッセージ例**:
- S6: 「まだ全てのピースが揃っていません（3/5）」
- S7: 「あなたの答えと一致しません。一字一句、正確に思い出してください。」

**復帰可能**: ✅ 可能

---

#### 3. NETWORK_ERROR (通信エラー)

**発生条件**:
- API呼び出しの失敗
- タイムアウト
- サーバーエラー

**動作**:
- ErrorModalで通知
- リトライボタン表示
- 現在のシーンに留まる

**ユーザーメッセージ**:
> 通信エラーが発生しました。もう一度お試しください。

**復帰可能**: ✅ 可能（リトライ）

---

#### 4. INIT_ERROR (初期化エラー)

**発生条件**:
- GameFlowManagerの初期化失敗
- セッション作成失敗

**動作**:
- ErrorModalで通知
- 再読み込みボタン表示

**ユーザーメッセージ**:
> ゲームの初期化に失敗しました

**復帰可能**: ⚠️ 部分的（再読み込みが必要）

---

### エラーモーダルの動作

**復帰可能なエラー**:
- モーダルを閉じられる
- 背景クリックで閉じる
- OKボタンのみ表示

**復帰不可能なエラー**:
- モーダルを閉じられない
- 2つのボタン表示:
  - 「再読み込み」
  - 「最初から」（S0へ）

---

## セッション管理

### セッションライフサイクル

```
1. セッション作成 (S0入場時)
   ↓
   session_id (UUID v4) 発行
   ↓
   LocalStorageに保存
   ↓
2. ゲームプレイ (S0〜S9)
   ↓
   各シーンで session_id を使用
   ↓
3. セッション終了
   ↓
   - 60分経過
   - または S9完了
   - または GAMEOVER
```

### セッション情報

```typescript
interface GameSession {
  sessionId: string;      // UUID v4
  createdAt: string;      // ISO 8601 UTC
  expiresAt: string;      // createdAt + 60分
  currentScene: string;   // 's0'〜's9', 'gameover'
}
```

### 保存場所

**クライアント側**:
```typescript
// LocalStorage
localStorage.setItem('game_session_id', sessionId);
localStorage.setItem('game_current_scene', currentScene);
```

**サーバー側**:
- MySQL sessions テーブル
- 60分後に自動削除（推奨）

---

## 開発者向け情報

### 遷移メソッドの使用

全シーンで `useGameFlow()` フックを使用:

```typescript
import { useGameFlow } from '@/hooks/useGameFlow';

export default function MyScenePage() {
  const { session, currentScene, transitionTo } = useGameFlow();

  const handleNext = async () => {
    // 自動的にバリデーション実行
    const success = await transitionTo('s2');

    if (!success) {
      // エラーはErrorModalで自動表示される
      console.log('遷移失敗');
    }
  };

  return (
    <button onClick={handleNext}>次へ</button>
  );
}
```

### バリデーションのスキップ

**通常は不要ですが**、デバッグ時など特殊な場合:

```typescript
// 第2引数に true を渡すとバリデーションをスキップ
await transitionTo('s7', true); // ⚠️ 本番では使用禁止
```

### 遷移可能性の事前チェック

```typescript
const { canTransitionTo } = useGameFlow();

// ボタンの有効/無効を動的に制御
const canProceed = await canTransitionTo('s7');

<Button disabled={!canProceed}>次へ</Button>
```

---

## まとめ

### 遷移パターン一覧

| 開始 | 終了 | 条件 | タイプ |
|------|------|------|--------|
| `/` | S0 | リンククリック | 通常 |
| S0 | S1 | 許可 + 同意 | バリデーション |
| S1 | S2 | 会話完了 | 通常 |
| S2 | S3 | 5箇所確認 | 通常 |
| S3 | S4 | ボタン押下 | 通常 |
| S4 | S5 | 全10問回答 | 通常 |
| S5 | S6 | ボタン押下 + API呼び出し | 通常 |
| S6 | S7 | 全5ピース取得 | バリデーション |
| S6 | GAMEOVER | 7分タイムアウト | タイマー |
| S7 | S8 | 回答完全一致 | バリデーション |
| S8 | S9 | 痕跡発見 | 通常 |
| S8 | GAMEOVER | 3分タイムアウト | タイマー |
| S9 | `/` | メッセージ保存後 | 通常 |
| GAMEOVER | S0 | 再挑戦ボタン | 通常 |
| GAMEOVER | `/` | トップへ戻る | 通常 |

### 重要なポイント

1. **S6が最重要**: 全ピース取得が必須
2. **S7が検証ポイント**: 完全一致が必須
3. **タイマーは2箇所**: S6（7分）、S8（3分）
4. **ブラウザバック無効**: 全シーンで戻る操作禁止
5. **セッション有効期限**: 60分

---

## 関連ファイル

### 実装ファイル
- 各シーンページ: `src/app/game/s[0-9]/page.tsx`
- ゲームフローマネージャー: `src/lib/game/GameFlowManager.ts`
- シーン遷移ルール: `src/lib/game/scene-transitions.ts`
- Context: `src/contexts/GameFlowContext.tsx`
- フック: `src/hooks/useGameFlow.ts`

### 定数ファイル
- ゲーム定数: `src/lib/game/constants.ts`
- タイマー設定: `TIMER_DURATIONS`, `TIMER_WARNING_THRESHOLDS`

### APIファイル
- API呼び出し: `src/lib/game/api.ts`

---

**最終更新**: 2025-11-08
**バージョン**: 1.0.0
