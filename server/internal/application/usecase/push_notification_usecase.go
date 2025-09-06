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

type SendPushRequest struct {
	UserID         *valueobject.UserID
	IdempotencyKey string
	Topic          string
	Urgency        model.Urgency
	TTLSeconds     int
	Payload        model.PushPayload
	ScheduleAt     *time.Time
}

type SendPushResponse struct {
	JobID   valueobject.JobID
	Success bool
	Message string
}

type SendBatchPushRequest struct {
	UserIDs        []valueobject.UserID
	Topic          string
	Urgency        model.Urgency
	TTLSeconds     int
	Payload        model.PushPayload
	ScheduleAt     *time.Time
	IdempotencyKey string
}

type SendBatchPushResponse struct {
	JobIDs  []valueobject.JobID
	Success bool
	Message string
}

type PushNotificationUseCase struct {
	jobRepo          repository.PushJobRepository
	subscriptionRepo repository.PushSubscriptionRepository
	pushService      *service.PushService
}

func NewPushNotificationUseCase(
	jobRepo repository.PushJobRepository,
	subscriptionRepo repository.PushSubscriptionRepository,
	pushService *service.PushService,
) *PushNotificationUseCase {
	return &PushNotificationUseCase{
		jobRepo:          jobRepo,
		subscriptionRepo: subscriptionRepo,
		pushService:      pushService,
	}
}

func (pnu *PushNotificationUseCase) SendPush(ctx context.Context, req SendPushRequest) (*SendPushResponse, error) {
	if req.IdempotencyKey != "" {
		existingJob, err := pnu.pushService.ValidateJobIdempotency(ctx, req.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("failed to validate idempotency: %w", err)
		}

		if existingJob != nil {
			return &SendPushResponse{
				JobID:   existingJob.ID(),
				Success: true,
				Message: "Job already exists (idempotent)",
			}, nil
		}
	}

	if req.UserID != nil {
		canReceive, err := pnu.pushService.CanUserReceivePush(ctx, *req.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to check if user can receive push: %w", err)
		}

		if !canReceive {
			return &SendPushResponse{
				Success: false,
				Message: "User has no valid push subscriptions",
			}, nil
		}
	}

	jobID, err := pnu.jobRepo.NextIdentity(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate job ID: %w", err)
	}

	if req.Urgency == "" {
		req.Urgency = model.UrgencyNormal
	}

	if req.TTLSeconds <= 0 {
		req.TTLSeconds = 86400
	}

	job, err := model.NewPushJob(
		jobID,
		req.IdempotencyKey,
		req.UserID,
		req.Topic,
		req.Urgency,
		req.TTLSeconds,
		req.Payload,
		req.ScheduleAt,
	)
	if err != nil {
		return &SendPushResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid job parameters: %v", err),
		}, nil
	}

	err = pnu.jobRepo.Save(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("failed to save push job: %w", err)
	}

	return &SendPushResponse{
		JobID:   jobID,
		Success: true,
		Message: "Push job created successfully",
	}, nil
}

func (pnu *PushNotificationUseCase) SendBatchPush(ctx context.Context, req SendBatchPushRequest) (*SendBatchPushResponse, error) {
	if req.IdempotencyKey != "" {
		existingJob, err := pnu.pushService.ValidateJobIdempotency(ctx, req.IdempotencyKey)
		if err != nil {
			return nil, fmt.Errorf("failed to validate idempotency: %w", err)
		}

		if existingJob != nil {
			return &SendBatchPushResponse{
				JobIDs:  []valueobject.JobID{existingJob.ID()},
				Success: true,
				Message: "Batch job already exists (idempotent)",
			}, nil
		}
	}

	var jobIDs []valueobject.JobID

	for _, userID := range req.UserIDs {
		canReceive, err := pnu.pushService.CanUserReceivePush(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check if user %s can receive push: %w", userID.String(), err)
		}

		if !canReceive {
			continue
		}

		jobID, err := pnu.jobRepo.NextIdentity(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate job ID: %w", err)
		}

		urgency := req.Urgency
		if urgency == "" {
			urgency = model.UrgencyNormal
		}

		ttlSeconds := req.TTLSeconds
		if ttlSeconds <= 0 {
			ttlSeconds = 86400
		}

		userIDCopy := userID
		job, err := model.NewPushJob(
			jobID,
			"",
			&userIDCopy,
			req.Topic,
			urgency,
			ttlSeconds,
			req.Payload,
			req.ScheduleAt,
		)
		if err != nil {
			return &SendBatchPushResponse{
				Success: false,
				Message: fmt.Sprintf("Invalid job parameters for user %s: %v", userID.String(), err),
			}, nil
		}

		err = pnu.jobRepo.Save(ctx, job)
		if err != nil {
			return nil, fmt.Errorf("failed to save push job for user %s: %w", userID.String(), err)
		}

		jobIDs = append(jobIDs, jobID)
	}

	return &SendBatchPushResponse{
		JobIDs:  jobIDs,
		Success: true,
		Message: fmt.Sprintf("Created %d push jobs successfully", len(jobIDs)),
	}, nil
}
