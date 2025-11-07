package dto

// SubscribeRequest は、プッシュ通知サブスクリプション登録リクエスト
type SubscribeRequest struct {
	SessionID string `json:"session_id"` // セッションID
	Endpoint  string `json:"endpoint"`   // プッシュエンドポイントURL
	Keys      struct {
		P256dh string `json:"p256dh"` // P256dh公開鍵（Base64）
		Auth   string `json:"auth"`   // 認証シークレット（Base64）
	} `json:"keys"`
}

// SubscribeResponse は、プッシュ通知サブスクリプション登録レスポンス
type SubscribeResponse struct {
	SubscriptionID string `json:"subscription_id"` // サブスクリプションID
	Message        string `json:"message"`         // メッセージ
}

// SendPushRequest は、プッシュ通知送信リクエスト
type SendPushRequest struct {
	Title   string `json:"title"`   // 通知タイトル
	Message string `json:"message"` // 通知メッセージ
}

// SendPushResponse は、プッシュ通知送信レスポンス
type SendPushResponse struct {
	Success bool   `json:"success"` // 成功フラグ
	Message string `json:"message"` // メッセージ
}

// VAPIDPublicKeyResponse は、VAPID公開鍵取得レスポンス
type VAPIDPublicKeyResponse struct {
	PublicKey string `json:"public_key"` // VAPID公開鍵（Base64 URL-safe）
}
