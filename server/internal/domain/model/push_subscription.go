package model

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type PushSubscription struct {
	id             valueobject.SubscriptionID
	userID         *valueobject.UserID 
	endpoint       valueobject.PushEndpoint
	keys           valueobject.PushKeys
	userAgent      string
	expirationTime *time.Time
	isValid        bool
	createdAt      time.Time
	updatedAt      time.Time
}

func NewPushSubscription(
	id valueobject.SubscriptionID,
	userID *valueobject.UserID,
	endpoint valueobject.PushEndpoint,
	keys valueobject.PushKeys,
	userAgent string,
	expirationTime *time.Time,
) *PushSubscription {
	now := time.Now()
	return &PushSubscription{
		id:             id,
		userID:         userID,
		endpoint:       endpoint,
		keys:           keys,
		userAgent:      userAgent,
		expirationTime: expirationTime,
		isValid:        true,
		createdAt:      now,
		updatedAt:      now,
	}
}

func ReconstructPushSubscription(
	id valueobject.SubscriptionID,
	userID *valueobject.UserID,
	endpoint valueobject.PushEndpoint,
	keys valueobject.PushKeys,
	userAgent string,
	expirationTime *time.Time,
	isValid bool,
	createdAt, updatedAt time.Time,
) *PushSubscription {
	return &PushSubscription{
		id:             id,
		userID:         userID,
		endpoint:       endpoint,
		keys:           keys,
		userAgent:      userAgent,
		expirationTime: expirationTime,
		isValid:        isValid,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

func (ps *PushSubscription) ID() valueobject.SubscriptionID {
	return ps.id
}

func (ps *PushSubscription) UserID() *valueobject.UserID {
	return ps.userID
}

func (ps *PushSubscription) Endpoint() valueobject.PushEndpoint {
	return ps.endpoint
}

func (ps *PushSubscription) Keys() valueobject.PushKeys {
	return ps.keys
}

func (ps *PushSubscription) UserAgent() string {
	return ps.userAgent
}

func (ps *PushSubscription) ExpirationTime() *time.Time {
	return ps.expirationTime
}

func (ps *PushSubscription) IsValid() bool {
	return ps.isValid
}

func (ps *PushSubscription) CreatedAt() time.Time {
	return ps.createdAt
}

func (ps *PushSubscription) UpdatedAt() time.Time {
	return ps.updatedAt
}

func (ps *PushSubscription) MarkAsInvalid() {
	ps.isValid = false
	ps.updatedAt = time.Now()
}

func (ps *PushSubscription) IsExpired() bool {
	if ps.expirationTime == nil {
		return false
	}
	return time.Now().After(*ps.expirationTime)
}

func (ps *PushSubscription) UpdateKeys(keys valueobject.PushKeys) {
	ps.keys = keys
	ps.updatedAt = time.Now()
}

func (ps *PushSubscription) UpdateUserAgent(userAgent string) {
	ps.userAgent = userAgent
	ps.updatedAt = time.Now()
}