package handler

import (
	"encoding/json"
	"net/http"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/dto"
)

// AnswerHandler は、S4回答管理のHTTPハンドラー
type AnswerHandler struct {
	s4AnswerUseCase *usecase.S4AnswerUseCase
}

// NewAnswerHandler は、新しいAnswerHandlerを作成する
func NewAnswerHandler(s4AnswerUseCase *usecase.S4AnswerUseCase) *AnswerHandler {
	return &AnswerHandler{
		s4AnswerUseCase: s4AnswerUseCase,
	}
}

// SaveAnswer は、S4の回答を保存する
// POST /api/session/{session_id}/answers
func (h *AnswerHandler) SaveAnswer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionIDStr := r.PathValue("session_id")

	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	var req dto.AnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "INVALID_REQUEST",
			"message": "invalid request body",
		})
		return
	}

	questionID, err := valueobject.NewQuestionID(req.QuestionID)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := h.s4AnswerUseCase.SaveAnswer(ctx, sessionID, questionID, req.AnswerText); err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"message": "answer saved successfully"})
}

// GetAnswers は、セッションの全回答を取得する
// GET /api/session/{session_id}/answers
func (h *AnswerHandler) GetAnswers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionIDStr := r.PathValue("session_id")

	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	answers, err := h.s4AnswerUseCase.GetAnswers(ctx, sessionID)
	if err != nil {
		handleError(w, err)
		return
	}

	responses := dto.ToAnswerResponses(answers)
	respondJSON(w, http.StatusOK, responses)
}
