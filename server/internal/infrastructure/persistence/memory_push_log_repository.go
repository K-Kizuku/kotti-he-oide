package persistence

import (
	"context"
	"sync"
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type MemoryPushLogRepository struct {
	mu     sync.RWMutex
	logs   map[int64]*model.PushLog
	nextID int64
}

func NewMemoryPushLogRepository() *MemoryPushLogRepository {
	return &MemoryPushLogRepository{
		logs:   make(map[int64]*model.PushLog),
		nextID: 1,
	}
}

func (r *MemoryPushLogRepository) Save(ctx context.Context, log *model.PushLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logs[log.ID()] = log
	return nil
}

func (r *MemoryPushLogRepository) FindByJobID(ctx context.Context, jobID valueobject.JobID) ([]*model.PushLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.PushLog
	for _, log := range r.logs {
		if log.JobID() != nil && log.JobID().Equals(jobID) {
			result = append(result, log)
		}
	}
	return result, nil
}

func (r *MemoryPushLogRepository) FindBySubscriptionID(ctx context.Context, subscriptionID valueobject.SubscriptionID) ([]*model.PushLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.PushLog
	for _, log := range r.logs {
		if log.SubscriptionID() != nil && log.SubscriptionID().Equals(subscriptionID) {
			result = append(result, log)
		}
	}
	return result, nil
}

func (r *MemoryPushLogRepository) DeleteOldLogs(ctx context.Context, olderThanDays int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().Add(-time.Duration(olderThanDays) * 24 * time.Hour)
	for id, log := range r.logs {
		if log.CreatedAt().Before(cutoff) {
			delete(r.logs, id)
		}
	}
	return nil
}

func (r *MemoryPushLogRepository) CountSuccessByJobID(ctx context.Context, jobID valueobject.JobID) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, log := range r.logs {
		if log.JobID() != nil && log.JobID().Equals(jobID) && log.IsSuccess() {
			count++
		}
	}
	return count, nil
}

func (r *MemoryPushLogRepository) CountFailuresByJobID(ctx context.Context, jobID valueobject.JobID) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, log := range r.logs {
		if log.JobID() != nil && log.JobID().Equals(jobID) && !log.IsSuccess() && log.ErrorMessage() != "" {
			count++
		}
	}
	return count, nil
}

func (r *MemoryPushLogRepository) NextIdentity(ctx context.Context) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++
	return id, nil
}