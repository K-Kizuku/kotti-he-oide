package repository

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// PushSubscriptionRepository は、プッシュ通知サブスクリプションの永続化を担当するリポジトリインターフェース
type PushSubscriptionRepository interface {
	// Save は、サブスクリプションを保存する
	Save(ctx context.Context, subscription *model.PushSubscription) error

	// FindByID は、サブスクリプションIDでサブスクリプションを検索する
	FindByID(ctx context.Context, subscriptionID valueobject.SubscriptionID) (*model.PushSubscription, error)

	// FindBySessionID は、セッションIDに紐づくアクティブなサブスクリプションを検索する
	FindBySessionID(ctx context.Context, sessionID valueobject.SessionID) (*model.PushSubscription, error)

	// Update は、サブスクリプションを更新する
	Update(ctx context.Context, subscription *model.PushSubscription) error

	// Delete は、サブスクリプションを削除する
	Delete(ctx context.Context, subscriptionID valueobject.SubscriptionID) error
}
