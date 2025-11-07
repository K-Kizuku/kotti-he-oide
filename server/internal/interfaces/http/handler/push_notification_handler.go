package handler

import (
	"encoding/json"
	"net/http"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/dto"
)

// PushNotificationHandler は、プッシュ通知のHTTPハンドラー
type PushNotificationHandler struct {
	pushNotificationUseCase usecase.PushNotificationUseCase
}

// NewPushNotificationHandler は、新しいPushNotificationHandlerを作成する
func NewPushNotificationHandler(pushNotificationUseCase usecase.PushNotificationUseCase) *PushNotificationHandler {
	return &PushNotificationHandler{
		pushNotificationUseCase: pushNotificationUseCase,
	}
}

// GetVAPIDPublicKey は、VAPID公開鍵を取得する
// GET /api/push/vapid-public-key
func (h *PushNotificationHandler) GetVAPIDPublicKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	publicKey, err := h.pushNotificationUseCase.GetVAPIDPublicKey(ctx)
	if err != nil {
		handleError(w, err)
		return
	}

	response := dto.VAPIDPublicKeyResponse{
		PublicKey: publicKey,
	}

	respondJSON(w, http.StatusOK, response)
}

// Subscribe は、プッシュ通知のサブスクリプションを登録する
// POST /api/push/subscribe
func (h *PushNotificationHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request body",
		})
		return
	}

	// バリデーション
	if req.SessionID == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "INVALID_INPUT",
			"message": "session_id is required",
		})
		return
	}

	if req.Endpoint == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "INVALID_INPUT",
			"message": "endpoint is required",
		})
		return
	}

	if req.Keys.P256dh == "" || req.Keys.Auth == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "INVALID_INPUT",
			"message": "keys (p256dh and auth) are required",
		})
		return
	}

	// SessionIDをValue Objectに変換
	sessionID, err := valueobject.NewSessionIDFromString(req.SessionID)
	if err != nil {
		handleError(w, err)
		return
	}

	// サブスクリプション登録
	subscription, err := h.pushNotificationUseCase.Subscribe(
		ctx,
		sessionID,
		req.Endpoint,
		req.Keys.P256dh,
		req.Keys.Auth,
	)

	if err != nil {
		handleError(w, err)
		return
	}

	response := dto.SubscribeResponse{
		SubscriptionID: subscription.SubscriptionID.String(),
		Message:        "Push notification subscription created successfully",
	}

	respondJSON(w, http.StatusCreated, response)
}

// Unsubscribe は、プッシュ通知のサブスクリプションを削除する
// DELETE /api/push/subscriptions/{subscription_id}
func (h *PushNotificationHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subscriptionIDStr := r.PathValue("subscription_id")

	subscriptionID, err := valueobject.NewSubscriptionIDFromString(subscriptionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := h.pushNotificationUseCase.Unsubscribe(ctx, subscriptionID); err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Push notification subscription deleted successfully",
	})
}

// SendPush は、特定のセッションに即時プッシュ通知を送信する
// POST /api/push/send/{session_id}
func (h *PushNotificationHandler) SendPush(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionIDStr := r.PathValue("session_id")

	// SessionIDをValue Objectに変換
	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	// リクエストボディをパース
	var req dto.SendPushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "INVALID_REQUEST",
			"message": "Invalid request body",
		})
		return
	}

	// バリデーション
	if req.Message == "" {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "INVALID_INPUT",
			"message": "message is required",
		})
		return
	}

	// プッシュ通知送信
	if err := h.pushNotificationUseCase.SendPushNotification(ctx, sessionID, req.Title, req.Message); err != nil {
		handleError(w, err)
		return
	}

	response := dto.SendPushResponse{
		Success: true,
		Message: "Push notification sent successfully",
	}

	respondJSON(w, http.StatusOK, response)
}
