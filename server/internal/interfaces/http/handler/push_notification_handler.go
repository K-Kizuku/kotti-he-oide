package handler

import (
	"encoding/json"
	"net/http"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/dto"
)

type PushNotificationHandler struct {
	notificationUseCase *usecase.PushNotificationUseCase
}

func NewPushNotificationHandler(notificationUseCase *usecase.PushNotificationUseCase) *PushNotificationHandler {
	return &PushNotificationHandler{
		notificationUseCase: notificationUseCase,
	}
}

func (pnh *PushNotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req dto.SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var userID *valueobject.UserID
	if req.UserID != nil && *req.UserID != "" {
		parsedUserID, err := valueobject.UserIDFromString(*req.UserID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		userID = &parsedUserID
	}

	urgency := model.Urgency(req.Urgency)
	if urgency == "" {
		urgency = model.UrgencyNormal
	}

	ttl := req.TTL
	if ttl <= 0 {
		ttl = 86400 // Default 24 hours
	}

	useCaseReq := usecase.SendPushRequest{
		UserID:         userID,
		IdempotencyKey: req.IdempotencyKey,
		Topic:          req.Topic,
		Urgency:        urgency,
		TTLSeconds:     ttl,
		Payload:        req.Payload,
		ScheduleAt:     req.ScheduleAt,
	}

	result, err := pnh.notificationUseCase.SendPush(r.Context(), useCaseReq)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := dto.SendNotificationResponse{
		JobID:   result.JobID.String(),
		Success: result.Success,
		Message: result.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	if result.Success {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(response)
}

func (pnh *PushNotificationHandler) SendBatchNotification(w http.ResponseWriter, r *http.Request) {
	var req dto.SendBatchNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var userIDs []valueobject.UserID
	for _, userIDStr := range req.UserIDs {
		userID, err := valueobject.UserIDFromString(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID: "+userIDStr, http.StatusBadRequest)
			return
		}
		userIDs = append(userIDs, userID)
	}

	urgency := model.Urgency(req.Urgency)
	if urgency == "" {
		urgency = model.UrgencyNormal
	}

	ttl := req.TTL
	if ttl <= 0 {
		ttl = 86400 // Default 24 hours
	}

	useCaseReq := usecase.SendBatchPushRequest{
		UserIDs:        userIDs,
		Topic:          req.Topic,
		Urgency:        urgency,
		TTLSeconds:     ttl,
		Payload:        req.Payload,
		ScheduleAt:     req.ScheduleAt,
		IdempotencyKey: req.IdempotencyKey,
	}

	result, err := pnh.notificationUseCase.SendBatchPush(r.Context(), useCaseReq)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	jobIDStrings := make([]string, len(result.JobIDs))
	for i, jobID := range result.JobIDs {
		jobIDStrings[i] = jobID.String()
	}

	response := dto.SendBatchNotificationResponse{
		JobIDs:  jobIDStrings,
		Success: result.Success,
		Message: result.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	if result.Success {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(response)
}