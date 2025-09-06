package repository

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type PushSubscriptionRepository interface {
	Save(ctx context.Context, subscription *model.PushSubscription) error
	FindByID(ctx context.Context, id valueobject.SubscriptionID) (*model.PushSubscription, error)
	FindByEndpoint(ctx context.Context, endpoint valueobject.PushEndpoint) (*model.PushSubscription, error)
	FindByUserID(ctx context.Context, userID valueobject.UserID) ([]*model.PushSubscription, error)
	FindValidSubscriptions(ctx context.Context) ([]*model.PushSubscription, error)
	FindValidSubscriptionsByUserID(ctx context.Context, userID valueobject.UserID) ([]*model.PushSubscription, error)
	MarkAsInvalid(ctx context.Context, id valueobject.SubscriptionID) error
	DeleteExpiredSubscriptions(ctx context.Context) error
	Delete(ctx context.Context, id valueobject.SubscriptionID) error
	NextIdentity(ctx context.Context) (valueobject.SubscriptionID, error)
}
