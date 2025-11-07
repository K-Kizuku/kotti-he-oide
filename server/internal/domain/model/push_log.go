package model

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// PushLog は、プッシュ通知送信ログを表すドメインモデル
type PushLog struct {
	ID             int64
	SubscriptionID valueobject.SubscriptionID
	SessionID      valueobject.SessionID
	Title          string
	Message        string
	Success        bool
	StatusCode     int
	ErrorMessage   string
	SentAt         time.Time
}

// NewPushLog は、新しいプッシュ通知送信ログを作成する（成功時）
func NewPushLog(
	subscriptionID valueobject.SubscriptionID,
	sessionID valueobject.SessionID,
	title string,
	message string,
	statusCode int,
) *PushLog {
	return &PushLog{
		SubscriptionID: subscriptionID,
		SessionID:      sessionID,
		Title:          title,
		Message:        message,
		Success:        true,
		StatusCode:     statusCode,
		ErrorMessage:   "",
		SentAt:         time.Now(),
	}
}

// NewPushLogWithError は、新しいプッシュ通知送信ログを作成する（失敗時）
func NewPushLogWithError(
	subscriptionID valueobject.SubscriptionID,
	sessionID valueobject.SessionID,
	title string,
	message string,
	statusCode int,
	errorMessage string,
) *PushLog {
	return &PushLog{
		SubscriptionID: subscriptionID,
		SessionID:      sessionID,
		Title:          title,
		Message:        message,
		Success:        false,
		StatusCode:     statusCode,
		ErrorMessage:   errorMessage,
		SentAt:         time.Now(),
	}
}

// IsSubscriptionGone は、サブスクリプションが無効になっているかを判定する（404/410エラー）
func (l *PushLog) IsSubscriptionGone() bool {
	return !l.Success && (l.StatusCode == 404 || l.StatusCode == 410)
}
