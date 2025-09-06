package model

import (
	"encoding/json"
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type PushLog struct {
	id              int64
	jobID           *valueobject.JobID
	subscriptionID  *valueobject.SubscriptionID
	responseStatus  *int
	responseHeaders map[string]string
	errorMessage    string
	createdAt       time.Time
}

func NewPushLog(
	id int64,
	jobID *valueobject.JobID,
	subscriptionID *valueobject.SubscriptionID,
	responseStatus *int,
	responseHeaders map[string]string,
	errorMessage string,
) *PushLog {
	return &PushLog{
		id:              id,
		jobID:           jobID,
		subscriptionID:  subscriptionID,
		responseStatus:  responseStatus,
		responseHeaders: responseHeaders,
		errorMessage:    errorMessage,
		createdAt:       time.Now(),
	}
}

func ReconstructPushLog(
	id int64,
	jobID *valueobject.JobID,
	subscriptionID *valueobject.SubscriptionID,
	responseStatus *int,
	responseHeaders map[string]string,
	errorMessage string,
	createdAt time.Time,
) *PushLog {
	return &PushLog{
		id:              id,
		jobID:           jobID,
		subscriptionID:  subscriptionID,
		responseStatus:  responseStatus,
		responseHeaders: responseHeaders,
		errorMessage:    errorMessage,
		createdAt:       createdAt,
	}
}

func (pl *PushLog) ID() int64 {
	return pl.id
}

func (pl *PushLog) JobID() *valueobject.JobID {
	return pl.jobID
}

func (pl *PushLog) SubscriptionID() *valueobject.SubscriptionID {
	return pl.subscriptionID
}

func (pl *PushLog) ResponseStatus() *int {
	return pl.responseStatus
}

func (pl *PushLog) ResponseHeaders() map[string]string {
	return pl.responseHeaders
}

func (pl *PushLog) ErrorMessage() string {
	return pl.errorMessage
}

func (pl *PushLog) CreatedAt() time.Time {
	return pl.createdAt
}

func (pl *PushLog) IsSuccess() bool {
	if pl.responseStatus == nil {
		return false
	}
	status := *pl.responseStatus
	return status >= 200 && status < 300
}

func (pl *PushLog) IsSubscriptionExpired() bool {
	if pl.responseStatus == nil {
		return false
	}
	status := *pl.responseStatus
	return status == 404 || status == 410
}

func (pl *PushLog) ResponseHeadersJSON() ([]byte, error) {
	if pl.responseHeaders == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(pl.responseHeaders)
}
