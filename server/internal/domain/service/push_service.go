package service

import (
	"context"
	"fmt"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type PushService struct {
	subscriptionRepo repository.PushSubscriptionRepository
	jobRepo          repository.PushJobRepository
}

func NewPushService(
	subscriptionRepo repository.PushSubscriptionRepository,
	jobRepo repository.PushJobRepository,
) *PushService {
	return &PushService{
		subscriptionRepo: subscriptionRepo,
		jobRepo:          jobRepo,
	}
}

func (ps *PushService) IsSubscriptionDuplicate(ctx context.Context, endpoint valueobject.PushEndpoint) (bool, error) {
	existing, err := ps.subscriptionRepo.FindByEndpoint(ctx, endpoint)
	if err != nil {
		return false, fmt.Errorf("failed to check duplicate subscription: %w", err)
	}
	return existing != nil, nil
}

func (ps *PushService) ValidateJobIdempotency(ctx context.Context, idempotencyKey string) (*model.PushJob, error) {
	if idempotencyKey == "" {
		return nil, nil
	}

	existing, err := ps.jobRepo.FindByIdempotencyKey(ctx, idempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("failed to check job idempotency: %w", err)
	}

	return existing, nil
}

func (ps *PushService) CanUserReceivePush(ctx context.Context, userID valueobject.UserID) (bool, error) {
	subscriptions, err := ps.subscriptionRepo.FindValidSubscriptionsByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user subscriptions: %w", err)
	}

	return len(subscriptions) > 0, nil
}

func (ps *PushService) CountActiveSubscriptions(ctx context.Context, userID valueobject.UserID) (int, error) {
	subscriptions, err := ps.subscriptionRepo.FindValidSubscriptionsByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user subscriptions: %w", err)
	}

	activeCount := 0
	for _, sub := range subscriptions {
		if !sub.IsExpired() {
			activeCount++
		}
	}

	return activeCount, nil
}

func (ps *PushService) CleanupExpiredSubscriptions(ctx context.Context) error {
	err := ps.subscriptionRepo.DeleteExpiredSubscriptions(ctx)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired subscriptions: %w", err)
	}
	return nil
}