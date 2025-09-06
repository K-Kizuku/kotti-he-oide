package persistence

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type MemoryPushSubscriptionRepository struct {
	mu            sync.RWMutex
	subscriptions map[valueobject.SubscriptionID]*model.PushSubscription
	nextID        int64
}

func NewMemoryPushSubscriptionRepository() *MemoryPushSubscriptionRepository {
	return &MemoryPushSubscriptionRepository{
		subscriptions: make(map[valueobject.SubscriptionID]*model.PushSubscription),
		nextID:        1,
	}
}

func (r *MemoryPushSubscriptionRepository) Save(ctx context.Context, subscription *model.PushSubscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.subscriptions[subscription.ID()] = subscription
	return nil
}

func (r *MemoryPushSubscriptionRepository) FindByID(ctx context.Context, id valueobject.SubscriptionID) (*model.PushSubscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	subscription, exists := r.subscriptions[id]
	if !exists {
		return nil, nil
	}
	return subscription, nil
}

func (r *MemoryPushSubscriptionRepository) FindByEndpoint(ctx context.Context, endpoint valueobject.PushEndpoint) (*model.PushSubscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, subscription := range r.subscriptions {
		if subscription.Endpoint().Equals(endpoint) {
			return subscription, nil
		}
	}
	return nil, nil
}

func (r *MemoryPushSubscriptionRepository) FindByUserID(ctx context.Context, userID valueobject.UserID) ([]*model.PushSubscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.PushSubscription
	for _, subscription := range r.subscriptions {
		if subscription.UserID() != nil && subscription.UserID().Equals(userID) {
			result = append(result, subscription)
		}
	}
	return result, nil
}

func (r *MemoryPushSubscriptionRepository) FindValidSubscriptions(ctx context.Context) ([]*model.PushSubscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.PushSubscription
	for _, subscription := range r.subscriptions {
		if subscription.IsValid() && !subscription.IsExpired() {
			result = append(result, subscription)
		}
	}
	return result, nil
}

func (r *MemoryPushSubscriptionRepository) FindValidSubscriptionsByUserID(ctx context.Context, userID valueobject.UserID) ([]*model.PushSubscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.PushSubscription
	for _, subscription := range r.subscriptions {
		if subscription.UserID() != nil && subscription.UserID().Equals(userID) && subscription.IsValid() && !subscription.IsExpired() {
			result = append(result, subscription)
		}
	}
	return result, nil
}

func (r *MemoryPushSubscriptionRepository) MarkAsInvalid(ctx context.Context, id valueobject.SubscriptionID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	subscription, exists := r.subscriptions[id]
	if !exists {
		return fmt.Errorf("subscription not found")
	}

	subscription.MarkAsInvalid()
	return nil
}

func (r *MemoryPushSubscriptionRepository) DeleteExpiredSubscriptions(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for id, subscription := range r.subscriptions {
		if subscription.ExpirationTime() != nil && now.After(*subscription.ExpirationTime()) {
			delete(r.subscriptions, id)
		}
	}
	return nil
}

func (r *MemoryPushSubscriptionRepository) Delete(ctx context.Context, id valueobject.SubscriptionID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.subscriptions, id)
	return nil
}

func (r *MemoryPushSubscriptionRepository) NextIdentity(ctx context.Context) (valueobject.SubscriptionID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++

	return valueobject.NewSubscriptionID(id)
}
