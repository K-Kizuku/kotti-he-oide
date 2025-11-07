package repository

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// PushLogRepository は、プッシュ通知送信ログの永続化を担当するリポジトリインターフェース
type PushLogRepository interface {
	// Save は、プッシュ通知送信ログを保存する
	Save(ctx context.Context, log *model.PushLog) error

	// FindBySessionID は、セッションIDに紐づくログを取得する
	FindBySessionID(ctx context.Context, sessionID valueobject.SessionID) ([]*model.PushLog, error)

	// FindBySubscriptionID は、サブスクリプションIDに紐づくログを取得する
	FindBySubscriptionID(ctx context.Context, subscriptionID valueobject.SubscriptionID) ([]*model.PushLog, error)
}
