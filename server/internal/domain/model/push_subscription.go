package model

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// PushSubscription は、プッシュ通知のサブスクリプション情報を表すドメインモデル
type PushSubscription struct {
	SubscriptionID valueobject.SubscriptionID
	SessionID      valueobject.SessionID
	Endpoint       valueobject.PushEndpoint
	Keys           valueobject.PushKeys
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewPushSubscription は、新しいプッシュ通知サブスクリプションを作成する
func NewPushSubscription(
	sessionID valueobject.SessionID,
	endpoint valueobject.PushEndpoint,
	keys valueobject.PushKeys,
) *PushSubscription {
	now := time.Now()
	return &PushSubscription{
		SubscriptionID: valueobject.NewSubscriptionID(),
		SessionID:      sessionID,
		Endpoint:       endpoint,
		Keys:           keys,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Deactivate は、サブスクリプションを無効化する（404/410エラー時など）
func (s *PushSubscription) Deactivate() error {
	if !s.IsActive {
		return errors.NewDomainError(
			errors.INVALID_STATE,
			"subscription is already inactive",
			nil,
		)
	}
	s.IsActive = false
	s.UpdatedAt = time.Now()
	return nil
}

// Reactivate は、サブスクリプションを再有効化する
func (s *PushSubscription) Reactivate() error {
	if s.IsActive {
		return errors.NewDomainError(
			errors.INVALID_STATE,
			"subscription is already active",
			nil,
		)
	}
	s.IsActive = true
	s.UpdatedAt = time.Now()
	return nil
}

// CanReceivePush は、プッシュ通知を受信可能かどうかを判定する
func (s *PushSubscription) CanReceivePush() bool {
	return s.IsActive
}
