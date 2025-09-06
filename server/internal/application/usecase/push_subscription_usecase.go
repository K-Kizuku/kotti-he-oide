package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type SubscribePushRequest struct {
	UserID         *valueobject.UserID
	Endpoint       string
	P256dhKey      string
	AuthKey        string
	UserAgent      string
	ExpirationTime *int64
}

type SubscribePushResponse struct {
	SubscriptionID valueobject.SubscriptionID
	Success        bool
	Message        string
}

type UnsubscribePushRequest struct {
	SubscriptionID valueobject.SubscriptionID
}

type UnsubscribePushResponse struct {
	Success bool
	Message string
}

type PushSubscriptionUseCase struct {
	subscriptionRepo repository.PushSubscriptionRepository
	pushService      *service.PushService
}

func NewPushSubscriptionUseCase(
	subscriptionRepo repository.PushSubscriptionRepository,
	pushService *service.PushService,
) *PushSubscriptionUseCase {
	return &PushSubscriptionUseCase{
		subscriptionRepo: subscriptionRepo,
		pushService:      pushService,
	}
}

func (psu *PushSubscriptionUseCase) Subscribe(ctx context.Context, req SubscribePushRequest) (*SubscribePushResponse, error) {
	endpoint, err := valueobject.NewPushEndpoint(req.Endpoint)
	if err != nil {
		return &SubscribePushResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid endpoint: %v", err),
		}, nil
	}

	p256dh, err := valueobject.NewP256dhKey(req.P256dhKey)
	if err != nil {
		return &SubscribePushResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid P256dh key: %v", err),
		}, nil
	}

	auth, err := valueobject.NewAuthKey(req.AuthKey)
	if err != nil {
		return &SubscribePushResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid auth key: %v", err),
		}, nil
	}

	isDuplicate, err := psu.pushService.IsSubscriptionDuplicate(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to check duplicate subscription: %w", err)
	}

	if isDuplicate {
		existing, err := psu.subscriptionRepo.FindByEndpoint(ctx, endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to find existing subscription: %w", err)
		}

		keys := valueobject.NewPushKeys(p256dh, auth)
		existing.UpdateKeys(keys)
		existing.UpdateUserAgent(req.UserAgent)

		err = psu.subscriptionRepo.Save(ctx, existing)
		if err != nil {
			return nil, fmt.Errorf("failed to update existing subscription: %w", err)
		}

		return &SubscribePushResponse{
			SubscriptionID: existing.ID(),
			Success:        true,
			Message:        "Subscription updated successfully",
		}, nil
	}

	id, err := psu.subscriptionRepo.NextIdentity(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate subscription ID: %w", err)
	}

	keys := valueobject.NewPushKeys(p256dh, auth)
	var expirationTime *time.Time
	if req.ExpirationTime != nil {
		t := time.Unix(*req.ExpirationTime/1000, 0)
		expirationTime = &t
	}

	subscription := model.NewPushSubscription(
		id,
		req.UserID,
		endpoint,
		keys,
		req.UserAgent,
		expirationTime,
	)

	err = psu.subscriptionRepo.Save(ctx, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to save subscription: %w", err)
	}

	return &SubscribePushResponse{
		SubscriptionID: id,
		Success:        true,
		Message:        "Subscription created successfully",
	}, nil
}

func (psu *PushSubscriptionUseCase) Unsubscribe(ctx context.Context, req UnsubscribePushRequest) (*UnsubscribePushResponse, error) {
	subscription, err := psu.subscriptionRepo.FindByID(ctx, req.SubscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find subscription: %w", err)
	}

	if subscription == nil {
		return &UnsubscribePushResponse{
			Success: false,
			Message: "Subscription not found",
		}, nil
	}

	subscription.MarkAsInvalid()
	err = psu.subscriptionRepo.Save(ctx, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to invalidate subscription: %w", err)
	}

	return &UnsubscribePushResponse{
		Success: true,
		Message: "Subscription removed successfully",
	}, nil
}
