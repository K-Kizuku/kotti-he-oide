package handler

import (
	"encoding/json"
	"net/http"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/dto"
)

type PushSubscriptionHandler struct {
	subscriptionUseCase *usecase.PushSubscriptionUseCase
}

func NewPushSubscriptionHandler(subscriptionUseCase *usecase.PushSubscriptionUseCase) *PushSubscriptionHandler {
	return &PushSubscriptionHandler{
		subscriptionUseCase: subscriptionUseCase,
	}
}

func (psh *PushSubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req dto.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Extract user ID from authentication context if needed
	var userID *valueobject.UserID

	useCaseReq := usecase.SubscribePushRequest{
		UserID:         userID,
		Endpoint:       req.Endpoint,
		P256dhKey:      req.Keys.P256dh,
		AuthKey:        req.Keys.Auth,
		UserAgent:      req.UserAgent,
		ExpirationTime: req.ExpirationTime,
	}

	result, err := psh.subscriptionUseCase.Subscribe(r.Context(), useCaseReq)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := dto.SubscribeResponse{
		ID:      result.SubscriptionID.String(),
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

func (psh *PushSubscriptionHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	subscriptionIDStr := r.PathValue("id")
	if subscriptionIDStr == "" {
		http.Error(w, "Subscription ID is required", http.StatusBadRequest)
		return
	}

	subscriptionID, err := valueobject.SubscriptionIDFromString(subscriptionIDStr)
	if err != nil {
		http.Error(w, "Invalid subscription ID", http.StatusBadRequest)
		return
	}

	useCaseReq := usecase.UnsubscribePushRequest{
		SubscriptionID: subscriptionID,
	}

	result, err := psh.subscriptionUseCase.Unsubscribe(r.Context(), useCaseReq)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := dto.UnsubscribeResponse{
		Success: result.Success,
		Message: result.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	if result.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(response)
}
