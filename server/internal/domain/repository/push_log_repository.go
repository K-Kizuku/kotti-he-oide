package repository

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type PushLogRepository interface {
	Save(ctx context.Context, log *model.PushLog) error
	FindByJobID(ctx context.Context, jobID valueobject.JobID) ([]*model.PushLog, error)
	FindBySubscriptionID(ctx context.Context, subscriptionID valueobject.SubscriptionID) ([]*model.PushLog, error)
	DeleteOldLogs(ctx context.Context, olderThanDays int) error
	CountSuccessByJobID(ctx context.Context, jobID valueobject.JobID) (int, error)
	CountFailuresByJobID(ctx context.Context, jobID valueobject.JobID) (int, error)
	NextIdentity(ctx context.Context) (int64, error)
}
