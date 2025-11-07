package handler

import (
	"encoding/json"
	"net/http"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/dto"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// SessionHandler は、セッション管理のHTTPハンドラー
type SessionHandler struct {
	sessionUseCase *usecase.SessionUseCase
}

// NewSessionHandler は、新しいSessionHandlerを作成する
func NewSessionHandler(sessionUseCase *usecase.SessionUseCase) *SessionHandler {
	return &SessionHandler{
		sessionUseCase: sessionUseCase,
	}
}

// CreateSession は、新しいセッションを作成する
// POST /api/session
func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, err := h.sessionUseCase.CreateSession(ctx)
	if err != nil {
		handleError(w, err)
		return
	}

	response := dto.ToSessionResponse(session)
	respondJSON(w, http.StatusCreated, response)
}

// GetSession は、セッション情報を取得する
// GET /api/session/{session_id}
func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionIDStr := r.PathValue("session_id")

	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	session, err := h.sessionUseCase.GetSession(ctx, sessionID)
	if err != nil {
		handleError(w, err)
		return
	}

	response := dto.ToSessionResponse(session)
	respondJSON(w, http.StatusOK, response)
}

// StartS6 は、S6を開始する（7分タイマー開始）
// POST /api/session/{session_id}/s6/start
func (h *SessionHandler) StartS6(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionIDStr := r.PathValue("session_id")

	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := h.sessionUseCase.StartS6(ctx, sessionID); err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "S6 started successfully"})
}

// handleError は、エラーをHTTPレスポンスに変換する
func handleError(w http.ResponseWriter, err error) {
	if domainErr, ok := err.(*errors.DomainError); ok {
		statusCode := http.StatusInternalServerError

		switch domainErr.Code {
		case errors.INVALID_INPUT, errors.INVALID_SESSION_ID:
			statusCode = http.StatusBadRequest
		case errors.SESSION_NOT_FOUND, errors.NOT_FOUND:
			statusCode = http.StatusNotFound
		case errors.SESSION_EXPIRED:
			statusCode = http.StatusGone
		case errors.S6_TIME_EXPIRED:
			statusCode = http.StatusRequestTimeout
		}

		respondJSON(w, statusCode, map[string]string{
			"error":   domainErr.Code,
			"message": domainErr.Message,
		})
		return
	}

	respondJSON(w, http.StatusInternalServerError, map[string]string{
		"error":   "INTERNAL_ERROR",
		"message": err.Error(),
	})
}

// respondJSON は、JSONレスポンスを返す
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
