package repository

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type PushJobRepository interface {
	Save(ctx context.Context, job *model.PushJob) error
	FindByID(ctx context.Context, id valueobject.JobID) (*model.PushJob, error)
	FindByIdempotencyKey(ctx context.Context, key string) (*model.PushJob, error)
	FindPendingJobs(ctx context.Context, limit int) ([]*model.PushJob, error)
	FindReadyToSendJobs(ctx context.Context, limit int) ([]*model.PushJob, error)
	FindFailedJobsForRetry(ctx context.Context, maxRetries int, limit int) ([]*model.PushJob, error)
	UpdateStatus(ctx context.Context, id valueobject.JobID, status model.JobStatus, lastError string) error
	IncrementRetryCount(ctx context.Context, id valueobject.JobID) error
	Delete(ctx context.Context, id valueobject.JobID) error
	DeleteOldCompletedJobs(ctx context.Context, olderThan int) error
	NextIdentity(ctx context.Context) (valueobject.JobID, error)
}