package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/dto"
)

// MessageHandler は、プレイヤーメッセージ管理のHTTPハンドラー
type MessageHandler struct {
	messageUseCase *usecase.MessageUseCase
}

// NewMessageHandler は、新しいMessageHandlerを作成する
func NewMessageHandler(messageUseCase *usecase.MessageUseCase) *MessageHandler {
	return &MessageHandler{
		messageUseCase: messageUseCase,
	}
}

// SaveMessage は、プレイヤーメッセージを保存する
// POST /api/session/{session_id}/message
func (h *MessageHandler) SaveMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionIDStr := r.PathValue("session_id")

	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	var req dto.MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "INVALID_REQUEST",
			"message": "invalid request body",
		})
		return
	}

	placeID, err := valueobject.NewPlaceID(req.PlaceID)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := h.messageUseCase.SaveMessage(ctx, sessionID, placeID, req.MessageText); err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"message": "message saved successfully"})
}

// GetMessages は、メッセージ一覧を取得する
// GET /api/messages
// GET /api/messages?place_id={place_id}&limit={limit}
func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// クエリパラメータを取得
	placeIDStr := r.URL.Query().Get("place_id")
	limitStr := r.URL.Query().Get("limit")

	// limitのデフォルト値は100
	limit := 100
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// place_idが指定されている場合は、その場所のメッセージのみを取得
	if placeIDStr != "" {
		placeID, err := valueobject.NewPlaceID(placeIDStr)
		if err != nil {
			handleError(w, err)
			return
		}

		messages, err := h.messageUseCase.GetMessagesByPlace(ctx, placeID, limit)
		if err != nil {
			handleError(w, err)
			return
		}

		responses := dto.ToMessageResponses(messages)
		respondJSON(w, http.StatusOK, responses)
		return
	}

	// place_idが指定されていない場合は、全メッセージを取得
	messages, err := h.messageUseCase.GetAllMessages(ctx, limit)
	if err != nil {
		handleError(w, err)
		return
	}

	responses := dto.ToMessageResponses(messages)
	respondJSON(w, http.StatusOK, responses)
}
