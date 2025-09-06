package dto

import "time"

type SubscribeRequest struct {
	Endpoint       string            `json:"endpoint" validate:"required,url"`
	Keys           PushKeys          `json:"keys" validate:"required"`
	UserAgent      string            `json:"ua,omitempty"`
	ExpirationTime *int64            `json:"expirationTime,omitempty"`
}

type PushKeys struct {
	P256dh string `json:"p256dh" validate:"required"`
	Auth   string `json:"auth" validate:"required"`
}

type SubscribeResponse struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type UnsubscribeRequest struct {
	ID string `json:"id" validate:"required"`
}

type UnsubscribeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SendNotificationRequest struct {
	UserID         *string                `json:"userId,omitempty"`
	IdempotencyKey string                 `json:"idempotencyKey,omitempty"`
	Topic          string                 `json:"topic,omitempty"`
	Urgency        string                 `json:"urgency,omitempty"`
	TTL            int                    `json:"ttl,omitempty"`
	Payload        map[string]interface{} `json:"payload" validate:"required"`
	ScheduleAt     *time.Time             `json:"scheduleAt,omitempty"`
}

type SendNotificationResponse struct {
	JobID   string `json:"jobId"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SendBatchNotificationRequest struct {
	UserIDs        []string               `json:"userIds" validate:"required,min=1"`
	Topic          string                 `json:"topic,omitempty"`
	Urgency        string                 `json:"urgency,omitempty"`
	TTL            int                    `json:"ttl,omitempty"`
	Payload        map[string]interface{} `json:"payload" validate:"required"`
	ScheduleAt     *time.Time             `json:"scheduleAt,omitempty"`
	IdempotencyKey string                 `json:"idempotencyKey,omitempty"`
}

type SendBatchNotificationResponse struct {
	JobIDs  []string `json:"jobIds"`
	Success bool     `json:"success"`
	Message string   `json:"message"`
}

type VAPIDPublicKeyResponse struct {
	PublicKey string `json:"publicKey"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

type ClickTrackingRequest struct {
	SubscriptionID string `json:"subscriptionId,omitempty"`
	JobID          string `json:"jobId,omitempty"`
	URL            string `json:"url,omitempty"`
	Timestamp      int64  `json:"timestamp"`
}