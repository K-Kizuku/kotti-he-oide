package persistence

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type MemoryPushJobRepository struct {
	mu     sync.RWMutex
	jobs   map[valueobject.JobID]*model.PushJob
	nextID int64
}

func NewMemoryPushJobRepository() *MemoryPushJobRepository {
	return &MemoryPushJobRepository{
		jobs:   make(map[valueobject.JobID]*model.PushJob),
		nextID: 1,
	}
}

func (r *MemoryPushJobRepository) Save(ctx context.Context, job *model.PushJob) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.jobs[job.ID()] = job
	return nil
}

func (r *MemoryPushJobRepository) FindByID(ctx context.Context, id valueobject.JobID) (*model.PushJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[id]
	if !exists {
		return nil, nil
	}
	return job, nil
}

func (r *MemoryPushJobRepository) FindByIdempotencyKey(ctx context.Context, key string) (*model.PushJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, job := range r.jobs {
		if job.IdempotencyKey() == key && key != "" {
			return job, nil
		}
	}
	return nil, nil
}

func (r *MemoryPushJobRepository) FindPendingJobs(ctx context.Context, limit int) ([]*model.PushJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.PushJob
	count := 0
	for _, job := range r.jobs {
		if job.Status() == model.JobStatusPending && count < limit {
			result = append(result, job)
			count++
		}
	}
	return result, nil
}

func (r *MemoryPushJobRepository) FindReadyToSendJobs(ctx context.Context, limit int) ([]*model.PushJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.PushJob
	count := 0
	for _, job := range r.jobs {
		if job.IsReadyToSend() && count < limit {
			result = append(result, job)
			count++
		}
	}
	return result, nil
}

func (r *MemoryPushJobRepository) FindFailedJobsForRetry(ctx context.Context, maxRetries int, limit int) ([]*model.PushJob, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.PushJob
	count := 0
	for _, job := range r.jobs {
		if job.ShouldRetry(maxRetries) && count < limit {
			result = append(result, job)
			count++
		}
	}
	return result, nil
}

func (r *MemoryPushJobRepository) UpdateStatus(ctx context.Context, id valueobject.JobID, status model.JobStatus, lastError string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[id]
	if !exists {
		return fmt.Errorf("job not found")
	}

	switch status {
	case model.JobStatusSending:
		job.MarkAsSending()
	case model.JobStatusSucceeded:
		job.MarkAsSucceeded()
	case model.JobStatusFailed:
		job.MarkAsFailed(lastError)
	case model.JobStatusCancelled:
		job.MarkAsCancelled()
	}

	return nil
}

func (r *MemoryPushJobRepository) IncrementRetryCount(ctx context.Context, id valueobject.JobID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	job, exists := r.jobs[id]
	if !exists {
		return fmt.Errorf("job not found")
	}

	job.MarkAsFailed(job.LastError())
	return nil
}

func (r *MemoryPushJobRepository) Delete(ctx context.Context, id valueobject.JobID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.jobs, id)
	return nil
}

func (r *MemoryPushJobRepository) DeleteOldCompletedJobs(ctx context.Context, olderThan int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().Add(-time.Duration(olderThan) * 24 * time.Hour)
	for id, job := range r.jobs {
		if (job.Status() == model.JobStatusSucceeded || job.Status() == model.JobStatusFailed) &&
			job.UpdatedAt().Before(cutoff) {
			delete(r.jobs, id)
		}
	}
	return nil
}

func (r *MemoryPushJobRepository) NextIdentity(ctx context.Context) (valueobject.JobID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++

	return valueobject.NewJobID(id)
}