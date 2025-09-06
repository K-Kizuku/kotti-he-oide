package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusSending   JobStatus = "sending"
	JobStatusSucceeded JobStatus = "succeeded"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

type Urgency string

const (
	UrgencyVeryLow Urgency = "very-low"
	UrgencyLow     Urgency = "low"
	UrgencyNormal  Urgency = "normal"
	UrgencyHigh    Urgency = "high"
)

func (u Urgency) IsValid() bool {
	switch u {
	case UrgencyVeryLow, UrgencyLow, UrgencyNormal, UrgencyHigh:
		return true
	default:
		return false
	}
}

type PushPayload map[string]interface{}

func (p PushPayload) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

type PushJob struct {
	id              valueobject.JobID
	idempotencyKey  string
	userID          *valueobject.UserID
	topic           string
	urgency         Urgency
	ttlSeconds      int
	payload         PushPayload
	scheduleAt      *time.Time
	status          JobStatus
	retryCount      int
	lastError       string
	createdAt       time.Time
	updatedAt       time.Time
}

func NewPushJob(
	id valueobject.JobID,
	idempotencyKey string,
	userID *valueobject.UserID,
	topic string,
	urgency Urgency,
	ttlSeconds int,
	payload PushPayload,
	scheduleAt *time.Time,
) (*PushJob, error) {
	if !urgency.IsValid() {
		return nil, fmt.Errorf("invalid urgency: %s", urgency)
	}

	if ttlSeconds < 0 {
		return nil, fmt.Errorf("TTL seconds must be non-negative")
	}

	now := time.Now()
	return &PushJob{
		id:              id,
		idempotencyKey:  idempotencyKey,
		userID:          userID,
		topic:           topic,
		urgency:         urgency,
		ttlSeconds:      ttlSeconds,
		payload:         payload,
		scheduleAt:      scheduleAt,
		status:          JobStatusPending,
		retryCount:      0,
		lastError:       "",
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

func ReconstructPushJob(
	id valueobject.JobID,
	idempotencyKey string,
	userID *valueobject.UserID,
	topic string,
	urgency Urgency,
	ttlSeconds int,
	payload PushPayload,
	scheduleAt *time.Time,
	status JobStatus,
	retryCount int,
	lastError string,
	createdAt, updatedAt time.Time,
) *PushJob {
	return &PushJob{
		id:              id,
		idempotencyKey:  idempotencyKey,
		userID:          userID,
		topic:           topic,
		urgency:         urgency,
		ttlSeconds:      ttlSeconds,
		payload:         payload,
		scheduleAt:      scheduleAt,
		status:          status,
		retryCount:      retryCount,
		lastError:       lastError,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

func (pj *PushJob) ID() valueobject.JobID {
	return pj.id
}

func (pj *PushJob) IdempotencyKey() string {
	return pj.idempotencyKey
}

func (pj *PushJob) UserID() *valueobject.UserID {
	return pj.userID
}

func (pj *PushJob) Topic() string {
	return pj.topic
}

func (pj *PushJob) Urgency() Urgency {
	return pj.urgency
}

func (pj *PushJob) TTLSeconds() int {
	return pj.ttlSeconds
}

func (pj *PushJob) Payload() PushPayload {
	return pj.payload
}

func (pj *PushJob) ScheduleAt() *time.Time {
	return pj.scheduleAt
}

func (pj *PushJob) Status() JobStatus {
	return pj.status
}

func (pj *PushJob) RetryCount() int {
	return pj.retryCount
}

func (pj *PushJob) LastError() string {
	return pj.lastError
}

func (pj *PushJob) CreatedAt() time.Time {
	return pj.createdAt
}

func (pj *PushJob) UpdatedAt() time.Time {
	return pj.updatedAt
}

func (pj *PushJob) MarkAsSending() {
	pj.status = JobStatusSending
	pj.updatedAt = time.Now()
}

func (pj *PushJob) MarkAsSucceeded() {
	pj.status = JobStatusSucceeded
	pj.lastError = ""
	pj.updatedAt = time.Now()
}

func (pj *PushJob) MarkAsFailed(error string) {
	pj.status = JobStatusFailed
	pj.lastError = error
	pj.retryCount++
	pj.updatedAt = time.Now()
}

func (pj *PushJob) MarkAsCancelled() {
	pj.status = JobStatusCancelled
	pj.updatedAt = time.Now()
}

func (pj *PushJob) IsReadyToSend() bool {
	if pj.status != JobStatusPending {
		return false
	}

	if pj.scheduleAt != nil && time.Now().Before(*pj.scheduleAt) {
		return false
	}

	return true
}

func (pj *PushJob) ShouldRetry(maxRetries int) bool {
	return pj.status == JobStatusFailed && pj.retryCount < maxRetries
}