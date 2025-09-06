# 1. 目的と要件整理

## 1.1 目的

- ブラウザの Push サービス（FCM/Firefox Autopush/WebKit Push 等）に対し、**Web Push プロトコル**で通知を送るアプリサーバを構築する。メッセージは RFC 8291 に従って暗号化、VAPID でサーバ識別。([datatracker.ietf.org][1])

## 1.2 必須要件

- HTTPS（サーバ ⇄Push サービス間も TLS 推奨）([tech-invite.com][3])
- VAPID 鍵（公開/秘密）管理とローテーション ([datatracker.ietf.org][4])
- 失効エンドポイント(404/410)のクリーンアップ ([tech-invite.com][3])
- 送信の TTL / Urgency / Topic（折り畳み）ヘッダ対応 ([web.dev][5], [Go Packages][6])
- サブスクリプションの保存・重複防止・失効監視（自動削除）

## 1.3 推奨要件

- ユーザー別・トピック別の配信設定
- 送信ジョブの非同期化（キュー/ワーカー）
- 監視（メトリクス/ログ/トレース）
- AB テスト・テンプレート化・言語/タイムゾーン差分

---

# 2. アーキテクチャ（推奨）

- **API サーバ（Go）**：REST（chi など）。購読登録/解除、公開鍵配布、配信要求受付、クリック計測。
- **DB（PostgreSQL）**：サブスクリプション・ユーザー・配信ジョブ・配信ログ。
- **キュー（例：GCP Pub/Sub or Redis）**：送信ワーカーのファンアウト/リトライ。
- **送信ワーカー（Go）**：`webpush-go` で Push サービスへ POST、結果を記録。([GitHub][7], [Go Packages][6])
- **シークレット管理**：VAPID 鍵（GCP Secret Manager 等）。
- **スケジューラ**：バッチ配信/定期クリーンアップ。

> Web Push の暗号化/プロトコル/VAPID の仕様根拠：RFC 8030/8291/8292、Push API/PushSubscription（MDN）。([datatracker.ietf.org][1], [MDN ウェブドキュメント][2])

---

# 3. データモデル（PostgreSQL）

```sql
-- users（任意: 自前認証を想定）
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email TEXT UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- push_subscriptions: ブラウザ購読
CREATE TABLE push_subscriptions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  endpoint TEXT NOT NULL UNIQUE,            -- 機微扱い（秘匿推奨）:contentReference[oaicite:8]{index=8}
  p256dh   TEXT NOT NULL,
  auth     TEXT NOT NULL,
  ua       TEXT,                            -- User-Agent 等
  expiration_time TIMESTAMPTZ,              -- MDNのexpire情報を保存すると良い:contentReference[oaicite:9]{index=9}
  is_valid BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 通知テンプレート（任意）
CREATE TABLE notification_templates (
  id BIGSERIAL PRIMARY KEY,
  key TEXT UNIQUE NOT NULL,
  title TEXT NOT NULL,
  body  TEXT NOT NULL,
  url   TEXT,
  icon  TEXT,
  data  JSONB,                              -- 任意の追加ペイロード
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ユーザーの配信設定
CREATE TABLE notification_prefs (
  user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  topics JSONB NOT NULL DEFAULT '{}'::jsonb, -- 例: {"news": true, "promo": false}
  quiet_hours JSONB                           -- サイレント時間帯など
);

-- 送信ジョブ（非同期）
CREATE TYPE job_status AS ENUM ('pending','sending','succeeded','failed','cancelled');
CREATE TABLE push_jobs (
  id BIGSERIAL PRIMARY KEY,
  idempotency_key TEXT UNIQUE,              -- 二重送信防止
  template_key TEXT REFERENCES notification_templates(key),
  user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  topic TEXT,                               -- Web PushのTopicヘッダ用（連続置換）:contentReference[oaicite:10]{index=10}
  urgency TEXT CHECK (urgency IN ('very-low','low','normal','high')), -- RFCの定義に準拠 :contentReference[oaicite:11]{index=11}
  ttl_seconds INT CHECK (ttl_seconds >= 0),
  payload JSONB NOT NULL,                   -- SWでJSON.parseされる前提
  schedule_at TIMESTAMPTZ,                  -- 未来配信に
  status job_status NOT NULL DEFAULT 'pending',
  retry_count INT NOT NULL DEFAULT 0,
  last_error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 配信ログ
CREATE TABLE push_logs (
  id BIGSERIAL PRIMARY KEY,
  job_id BIGINT REFERENCES push_jobs(id) ON DELETE SET NULL,
  subscription_id BIGINT REFERENCES push_subscriptions(id) ON DELETE SET NULL,
  response_status INT,
  response_headers JSONB,
  error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

# 4. 依存ライブラリ（Go）

- Router: `github.com/go-chi/chi/v5`（任意）
- DB: `github.com/jackc/pgx/v5`（または sqlc）
- Web Push: `github.com/SherClockHolmes/webpush-go`（VAPID/暗号化/Topic/TTL/Urgency をサポート）([Go Packages][6])
- ログ: `zap` など
- メトリクス: `prometheus/client_golang`

---

# 5. VAPID 鍵の生成・保管・ローテーション

## 5.1 生成（ワンショット）

```go
privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
```

公開鍵はフロントに配布、秘密鍵は Secret Manager 等で厳重保管。([Go Packages][6])

## 5.2 公開鍵 API

- `GET /push/vapid-public-key` → Base64URL 公開鍵を返す（SW/フロント購読に必要）。([MDN ウェブドキュメント][8])

## 5.3 ローテーション方針

- 「鍵バージョン」をメタに持ち、**新旧の公開鍵をしばらく同時配布**。購読更新後に古い秘密鍵を破棄。
- 送信側は「サブスクリプション取得時の鍵バージョン」を見て適合鍵で送る（段階的移行）。

---

# 6. REST API 設計（サーバ）

## 6.1 エンドポイント一覧

- `POST /api/push/subscribe`
  body: `PushSubscription`（`endpoint`, `keys.p256dh`, `keys.auth` など）＋任意メタ
  動作: 保存/更新（UPSERT）。返り値: 201。([MDN ウェブドキュメント][9])
- `DELETE /api/push/subscriptions/{id}`
  動作: 論理削除（`is_valid=false`）
- `GET /api/push/vapid-public-key`
  動作: 公開鍵を返す（Base64URL）
- `POST /api/push/send`（管理用・単発）
  body: `user_id` or `subscription_id` or `topic/broadcast` + `payload` + `ttl/urgency/topic`
  動作: push_jobs に enq
- `POST /api/push/send/batch`（管理用・条件配信）
- `POST /api/push/click`（通知クリック計測：SW から遷移時に ping）

> `PushSubscription` の構造や `endpoint` の意味は MDN。`endpoint` は秘匿推奨（乱用防止）。([MDN ウェブドキュメント][9])

## 6.2 ハンドラ例（購読登録）

```go
type SubscribeReq struct {
  Endpoint string `json:"endpoint" validate:"required,url"`
  Keys struct {
    P256dh string `json:"p256dh" validate:"required"`
    Auth   string `json:"auth"   validate:"required"`
  } `json:"keys"`
  UA string `json:"ua"`
  Expiration *time.Time `json:"expirationTime"`
}

func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
  var req SubscribeReq
  if err := json.NewDecoder(r.Body).Decode(&req); err != nil { http.Error(w, err.Error(), 400); return }
  // 認証済みユーザーIDを取り出す想定
  uid := authUserID(r)

  // UPSERT（endpoint は UNIQUE）
  id, err := h.store.UpsertSubscription(r.Context(), UpsertSubscriptionParams{
    UserID: uid, Endpoint: req.Endpoint, P256dh: req.Keys.P256dh, Auth: req.Keys.Auth,
    UA: req.UA, Expiration: req.Expiration,
  })
  if err != nil { http.Error(w, err.Error(), 500); return }

  w.WriteHeader(http.StatusCreated) // 201
  json.NewEncoder(w).Encode(map[string]any{"id": id})
}
```

---

# 7. 送信ワーカー（核心）

## 7.1 送信ライブラリとヘッダ

- `webpush-go` の `Options` で **TTL / Urgency / Topic** を指定可能。Topic は「未配信の同一トピックを置き換える」用途。([Go Packages][6])
- **Urgency の選択肢**：`very-low | low | normal | high`（デバイス状態に応じた配信ポリシー）。([web.dev][5])

## 7.2 送信処理（擬似コード）

```go
opts := &webpush.Options{
  Subscriber:      "mailto:ops@example.com",
  VAPIDPrivateKey: vapidPriv, // Secret Manager から取得
  TTL:             job.TTLSeconds,
  Urgency:         webpush.Urgency(job.Urgency), // "normal" など
  Topic:           job.Topic,                     // 折り畳み
}

payload, _ := json.Marshal(job.Payload) // SW側は event.data.json() で受ける

for _, sub := range subs { // 対象購読にファンアウト
  resp, err := webpush.SendNotificationWithContext(ctx, payload, &webpush.Subscription{
    Endpoint: sub.Endpoint,
    Keys: webpush.Keys{P256dh: sub.P256dh, Auth: sub.Auth},
  }, opts)

  if err != nil {
    recordError(job, sub, err)        // ネットワーク/暗号化エラー
    scheduleRetry(job, sub)           // バックオフ
    continue
  }
  defer resp.Body.Close()

  // 404/410 は購読失効 → is_valid=false にしてクリーンアップ
  if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
    markSubscriptionInvalid(sub)
  }

  // 2xx(200/201/202) は成功扱い（Push サービス側で受理）。:contentReference[oaicite:19]{index=19}
  recordResult(job, sub, resp.StatusCode, headersToJSON(resp.Header))
}
```

> 成功時の具体的ステータスは Push サービスにより `200/201/202` が返りうる。\*\*成功＝ユーザー端末に表示保証ではなく「Push サービスが受理」\*\*という意味なので、クリック計測やアプリ内既読で効果検証を行う。([web.dev][5], [pushpad.xyz][10])

## 7.3 リトライ戦略

- **指数バックオフ**（例：30s → 2m → 10m → 1h → 6h、上限 5 回）
- 恒久的エラー（4xx でバリデーション/認可失敗）は**即座に打ち切り**
- 一時的エラー（5xx/ネットワーク）は**再試行**
- 失効（404/410）は**購読無効化**（RFC 8030 に明記）。([tech-invite.com][3])

---

# 8. 配信制御（品質/到達率/体験）

## 8.1 TTL/Urgency/Topic の設計

- **TTL**：既読性の低い情報は短め（例：在庫/ライブ更新＝ 60〜300 秒）。ゆったり情報は長め（例：お知らせ＝ 1 日）。([web.dev][5])
- **Urgency**：通常 `normal`、本当に重要な場合のみ `high`（電池配慮）。([web.dev][5])
- **Topic**：同一トピックの未配信通知を**置き換え**（在庫数・スコア更新などに最適）。([Go Packages][6])

## 8.2 ペイロード設計（SW での表示）

- SW 側は `event.data?.json()` を想定し、`{title, body, url, icon, tag, data...}` を送る。
- 本文が長い/更新されやすい場合は **URL を渡し、開封後に API で取得**（軽量化）。

---

# 9. セキュリティ・プライバシー

- `endpoint` は第三者に知られると悪用されるため**秘匿情報として保存**（ログ赤伏せ）。([MDN ウェブドキュメント][11])
- VAPID 秘密鍵は Secret Manager、**権限分離**、**監査ログ**。
- API は認証必須（ユーザー紐づけがない匿名購読は、**レート制限**と **CAPTCHA** 検討）。
- 個人データ（地域/言語/クリック等）を扱うなら**プライバシーポリシー**と**同意管理**。
- 監査用に**誰に/何を/いつ/どの理由で**をログ化（PII はハッシュ化）。

---

# 10. 監視・可観測性

- **メトリクス**（Prometheus）

  - `push_requests_total{status,provider}`
  - `push_success_total` / `push_failed_total{code}`
  - `invalidated_subscriptions_total`
  - レイテンシ、キュー滞留長、リトライ回数

- **ログ**：構造化（zap）。`endpoint` はハッシュ化。
- **トレース**：送信ジョブ → ワーカー →Push サービス HTTP の span。

---

# 11. バックエンド実装（要点コード）

## 11.1 送信サービス

```go
type Sender struct {
  privKey string
  httpc   *http.Client
}

func NewSender(priv string) *Sender {
  return &Sender{privKey: priv, httpc: &http.Client{Timeout: 10 * time.Second}}
}

func (s *Sender) Send(ctx context.Context, sub SubscriptionRow, job PushJobRow) (*http.Response, error) {
  payload, _ := json.Marshal(job.Payload)
  opts := &webpush.Options{
    Subscriber:      "mailto:ops@example.com",
    VAPIDPrivateKey: s.privKey,
    TTL:             job.TTLSeconds,
    Urgency:         webpush.Urgency(job.Urgency),
    Topic:           job.Topic,
  }
  return webpush.SendNotificationWithContext(ctx, payload, &webpush.Subscription{
    Endpoint: sub.Endpoint,
    Keys: webpush.Keys{P256dh: sub.P256dh, Auth: sub.Auth},
  }, opts)
}
```

**webpush-go** の Options に `TTL/Urgency/Topic` が用意されていることは公式ドキュメントに記載。([Go Packages][6])

## 11.2 ワーカー（キュー消費）

- GCP **Pub/Sub Pull**（Cloud Run Jobs でも可）or Redis Stream。
- 1 メッセージ = 1 ジョブ or 1 ユーザー分のバッチ。

```go
for msg := range queue.Receive(ctx) {
  job := decodeJob(msg)
  subs := repo.ListValidSubscriptions(ctx, job.UserID)
  for _, sub := range subs {
    resp, err := sender.Send(ctx, sub, job)
    handleResult(job, sub, resp, err)
  }
  queue.Ack(msg)
}
```

---

# 12. 失効検知・クリーンアップ

- 送信結果の **404/410** は該当 `subscription` を `is_valid=false` に更新し、**後続の送信対象から除外**。RFC 8030 で「失効購読は 404 を返す」旨が規定。([tech-invite.com][3])
- バックグラウンドで `is_valid=false` を一定期間後に削除（個人データ最小化）。

---

# 13. 運用タスク

- **定期バッチ**

  - 期限切れ購読の削除
  - 古いジョブ/ログのアーカイブ
  - 未送信ジョブの再スケジュール

- **鍵ローテーション**

  - 新鍵を追加 → 公開鍵 API に両方返す → N 日後に旧鍵停止

- **負荷対策**

  - ワーカーの水平スケール（キュー滞留に応じた HPA）
  - Push サービス側の**レート制限**を考慮した送信速度制御

- **法令/同意**

  - オプトイン・オプトアウトの明示。**ユーザーによる解除 API** を用意。

---

# 14. iOS/Safari 実装に関する補足（サーバ目線）

- iOS/iPadOS の Web Push（16.4+）も**サーバは標準の Web Push プロトコル**で送るだけでよい（クライアント側は PWA 条件などの制約あり）。サーバ側は特別な APNs 資格情報は不要。クライアントの購読が得られれば通常通り送信できる。根拠：Push API/W3C、MDN。([W3C][12], [MDN ウェブドキュメント][2])

---

# 15. デプロイ（例：GCP Cloud Run）

1. **Artifact Registry** にビルド & Push
2. **Cloud Run** にデプロイ（最小 0 / 最大同時実行は API とワーカーで分離）
3. **Secret Manager**：`VAPID_PRIVATE_KEY` 注入
4. **Cloud SQL (Postgres)** 接続（pgx）
5. **Pub/Sub**（トピック `webpush.send` / サブスクリプション `worker`）
6. **Cloud Scheduler**：クリーンアップ/予約ジョブトリガ
7. **モニタリング**：Cloud Monitoring + Prometheus/Grafana

（あなたは Cloud Run/Artifact Registry/Workload Identity に慣れているので、この構成は親和性が高いはずです）

---

# 16. テスト

- **ユニット**：送信サービスを `httptest.Server` で差し替え、返すステータス（200/201/202/404/410/5xx）に応じて分岐を検証。
- **統合**：実ブラウザで購読を取得 → テスト環境から送信 →SW の `push` 受信を E2E で確認。
- **負荷**：大量購読へのファンアウト速度・Push サービス応答を測定（QPS/秒あたり送信上限の探索）。

---

# 17. 参考リンク（主要仕様/一次情報）

- **RFC 8030**（Web Push プロトコル・失効時の 404/410 など） ([datatracker.ietf.org][1], [tech-invite.com][3])
- **RFC 8291**（Web Push メッセージ暗号化） ([RFC Editor][13], [datatracker.ietf.org][14])
- **RFC 8292**（VAPID） ([datatracker.ietf.org][4], [RFC Editor][15])
- **Push API (MDN/W3C)**（購読/PushSubscription/subscribe） ([MDN ウェブドキュメント][2], [W3C][12])
- **web.dev: Web Push プロトコル**（Urgency/TTL/レスポンス確認） ([web.dev][5])
- **webpush-go（Go ライブラリ）**（Options: TTL/Topic/Urgency 等） ([Go Packages][6])
- **PushSubscription endpoint の秘匿**（MDN） ([MDN ウェブドキュメント][11])

---

## 付録 A：サンプル Dockerfile（マルチステージ）

```dockerfile
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server

FROM gcr.io/distroless/base-debian12
ENV PORT=8080
COPY --from=build /src/app /app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app"]
```

## 付録 B：通知ペイロードの例（サーバ →SW）

```json
{
  "title": "新着メッセージ",
  "body": "10件の更新があります（タップで確認）",
  "url": "/inbox",
  "icon": "/icons/icon-192.png",
  "tag": "inbox", // クライアント側の置き換え用
  "data": { "source": "batch-update" }
}
```

---

必要なら、**sqlc 用のクエリ定義**、**chi ルーター一式**、**Pub/Sub 実装**、**Prometheus メトリクス**まで丸ごと雛形を出します。どの部分からコード化しますか？（例：DB スキーマ →sqlc→REST→ ワーカー →Cloud Run 用 Terraform）

[1]: https://datatracker.ietf.org/doc/html/rfc8030?utm_source=chatgpt.com "RFC 8030 - Generic Event Delivery Using HTTP Push"
[2]: https://developer.mozilla.org/en-US/docs/Web/API/Push_API?utm_source=chatgpt.com "Push API - Web - MDN - Mozilla"
[3]: https://www.tech-invite.com/y80/tinv-ietf-rfc-8030-2.html?utm_source=chatgpt.com "RFC 8030: 2 of 2, p. 20 to 31"
[4]: https://datatracker.ietf.org/doc/html/rfc8292?utm_source=chatgpt.com "RFC 8292 - Voluntary Application Server Identification ..."
[5]: https://web.dev/articles/push-notifications-web-push-protocol?utm_source=chatgpt.com "The Web Push Protocol | Articles"
[6]: https://pkg.go.dev/github.com/sherclockholmes/webpush-go "webpush package - github.com/sherclockholmes/webpush-go - Go Packages"
[7]: https://github.com/SherClockHolmes/webpush-go?utm_source=chatgpt.com "SherClockHolmes/webpush-go: Web Push API Encryption ..."
[8]: https://developer.mozilla.org/en-US/docs/Web/API/PushManager/subscribe?utm_source=chatgpt.com "PushManager: subscribe() method - Web APIs - MDN - Mozilla"
[9]: https://developer.mozilla.org/en-US/docs/Web/API/PushSubscription?utm_source=chatgpt.com "PushSubscription - Web APIs - MDN - Mozilla"
[10]: https://pushpad.xyz/blog/web-push-errors-explained-with-http-status-codes?utm_source=chatgpt.com "Web Push errors explained (with HTTP status codes)"
[11]: https://developer.mozilla.org/en-US/docs/Web/API/PushSubscription/endpoint?utm_source=chatgpt.com "PushSubscription: endpoint property - Web APIs - MDN - Mozilla"
[12]: https://www.w3.org/TR/push-api/?utm_source=chatgpt.com "Push API"
[13]: https://www.rfc-editor.org/rfc/rfc8291.html?utm_source=chatgpt.com "RFC 8291: Message Encryption for Web Push"
[14]: https://datatracker.ietf.org/doc/html/rfc8291?utm_source=chatgpt.com "RFC 8291 - Message Encryption for Web Push"
[15]: https://www.rfc-editor.org/info/rfc8292?utm_source=chatgpt.com "Information on RFC 8292"
