package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
)

type PushSenderService struct {
	subscriptionRepo repository.PushSubscriptionRepository
	jobRepo          repository.PushJobRepository
	logRepo          repository.PushLogRepository
	vapidService     *service.VAPIDService
	httpClient       *http.Client
}

func NewPushSenderService(
	subscriptionRepo repository.PushSubscriptionRepository,
	jobRepo repository.PushJobRepository,
	logRepo repository.PushLogRepository,
	vapidService *service.VAPIDService,
) *PushSenderService {
	return &PushSenderService{
		subscriptionRepo: subscriptionRepo,
		jobRepo:          jobRepo,
		logRepo:          logRepo,
		vapidService:     vapidService,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (pss *PushSenderService) ProcessPendingJobs(ctx context.Context, batchSize int) error {
	jobs, err := pss.jobRepo.FindReadyToSendJobs(ctx, batchSize)
	if err != nil {
		return fmt.Errorf("failed to fetch ready jobs: %w", err)
	}

	for _, job := range jobs {
		err := pss.processJob(ctx, job)
		if err != nil {
			log.Printf("Failed to process job %s: %v", job.ID().String(), err)
		}
	}

	return nil
}

func (pss *PushSenderService) processJob(ctx context.Context, job *model.PushJob) error {
	job.MarkAsSending()
	err := pss.jobRepo.Save(ctx, job)
	if err != nil {
		return fmt.Errorf("failed to mark job as sending: %w", err)
	}

	var subscriptions []*model.PushSubscription

	if job.UserID() != nil {
		userSubs, err := pss.subscriptionRepo.FindValidSubscriptionsByUserID(ctx, *job.UserID())
		if err != nil {
			job.MarkAsFailed(fmt.Sprintf("Failed to get user subscriptions: %v", err))
			pss.jobRepo.Save(ctx, job)
			return fmt.Errorf("failed to get user subscriptions: %w", err)
		}
		subscriptions = userSubs
	} else {
		allSubs, err := pss.subscriptionRepo.FindValidSubscriptions(ctx)
		if err != nil {
			job.MarkAsFailed(fmt.Sprintf("Failed to get all subscriptions: %v", err))
			pss.jobRepo.Save(ctx, job)
			return fmt.Errorf("failed to get all subscriptions: %w", err)
		}
		subscriptions = allSubs
	}

	if len(subscriptions) == 0 {
		job.MarkAsSucceeded()
		pss.jobRepo.Save(ctx, job)
		return nil
	}

	successCount := 0
	failureCount := 0

	for _, subscription := range subscriptions {
		success, err := pss.sendToSubscription(ctx, job, subscription)
		if err != nil {
			log.Printf("Failed to send to subscription %s: %v", subscription.ID().String(), err)
			failureCount++
		} else if success {
			successCount++
		} else {
			failureCount++
		}
	}

	if successCount > 0 && failureCount == 0 {
		job.MarkAsSucceeded()
	} else if successCount == 0 && failureCount > 0 {
		job.MarkAsFailed(fmt.Sprintf("All %d deliveries failed", failureCount))
	} else {
		job.MarkAsSucceeded()
		log.Printf("Job %s completed with partial success: %d succeeded, %d failed",
			job.ID().String(), successCount, failureCount)
	}

	return pss.jobRepo.Save(ctx, job)
}

func (pss *PushSenderService) sendToSubscription(
	ctx context.Context,
	job *model.PushJob,
	subscription *model.PushSubscription,
) (bool, error) {
	payload, err := json.Marshal(job.Payload())
	if err != nil {
		return false, fmt.Errorf("failed to marshal payload: %w", err)
	}

	options := &webpush.Options{
		Subscriber:      "mailto:support@example.com",
		VAPIDPrivateKey: pss.vapidService.GetPrivateKey(),
		TTL:             job.TTLSeconds(),
		Urgency:         webpush.Urgency(job.Urgency()),
	}

	if job.Topic() != "" {
		options.Topic = job.Topic()
	}

	webpushSubscription := &webpush.Subscription{
		Endpoint: subscription.Endpoint().Value(),
		Keys: webpush.Keys{
			P256dh: subscription.Keys().P256dh().Value(),
			Auth:   subscription.Keys().Auth().Value(),
		},
	}

	resp, err := webpush.SendNotificationWithContext(ctx, payload, webpushSubscription, options)

	logID, _ := pss.logRepo.NextIdentity(ctx)

	if err != nil {
		jobID := job.ID()
		subscriptionID := subscription.ID()
		pushLog := model.NewPushLog(
			logID,
			&jobID,
			&subscriptionID,
			nil,
			nil,
			fmt.Sprintf("Send error: %v", err),
		)
		pss.logRepo.Save(ctx, pushLog)
		return false, err
	}

	defer resp.Body.Close()

	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}

	jobID := job.ID()
	subscriptionID := subscription.ID()
	pushLog := model.NewPushLog(
		logID,
		&jobID,
		&subscriptionID,
		&resp.StatusCode,
		responseHeaders,
		"",
	)

	pss.logRepo.Save(ctx, pushLog)

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
		subscription.MarkAsInvalid()
		pss.subscriptionRepo.Save(ctx, subscription)
		log.Printf("Marked subscription %s as invalid due to %d response",
			subscription.ID().String(), resp.StatusCode)
		return false, nil
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	}

	return false, fmt.Errorf("push service responded with status %d", resp.StatusCode)
}

func (pss *PushSenderService) ProcessRetries(ctx context.Context, maxRetries, batchSize int) error {
	jobs, err := pss.jobRepo.FindFailedJobsForRetry(ctx, maxRetries, batchSize)
	if err != nil {
		return fmt.Errorf("failed to fetch failed jobs for retry: %w", err)
	}

	for _, job := range jobs {
		backoffDelay := pss.calculateBackoffDelay(job.RetryCount())

		if time.Since(job.UpdatedAt()) < backoffDelay {
			continue
		}

		err := pss.processJob(ctx, job)
		if err != nil {
			log.Printf("Failed to retry job %s: %v", job.ID().String(), err)
		}
	}

	return nil
}

func (pss *PushSenderService) calculateBackoffDelay(retryCount int) time.Duration {
	baseDelay := 30 * time.Second
	maxDelay := 6 * time.Hour

	delay := baseDelay
	for i := 0; i < retryCount; i++ {
		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
			break
		}
	}

	return delay
}
