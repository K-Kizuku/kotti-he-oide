package handler

import (
	"encoding/json"
	"net/http"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/dto"
)

// S6Handler は、S6進捗管理のHTTPハンドラー
type S6Handler struct {
	s6UseCase *usecase.S6UseCase
}

// NewS6Handler は、新しいS6Handlerを作成する
func NewS6Handler(s6UseCase *usecase.S6UseCase) *S6Handler {
	return &S6Handler{
		s6UseCase: s6UseCase,
	}
}

// InitializeProgress は、S6の進捗を初期化する（5箇所分）
// POST /api/session/{session_id}/s6/initialize
func (h *S6Handler) InitializeProgress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionIDStr := r.PathValue("session_id")

	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	if err := h.s6UseCase.InitializeProgress(ctx, sessionID); err != nil {
		handleError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{"message": "progress initialized successfully"})
}

// VerifyLocation は、場所到達を記録する（Web選択のみ）
// POST /api/session/{session_id}/s6/verify-location
func (h *S6Handler) VerifyLocation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionIDStr := r.PathValue("session_id")

	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	var req dto.VerifyLocationRequest
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

	// Web選択として記録
	if err := h.s6UseCase.MarkPlaceVerified(ctx, sessionID, placeID, "manual"); err != nil {
		handleError(w, err)
		return
	}

	response := dto.VerifyLocationResponse{Verified: true}
	respondJSON(w, http.StatusOK, response)
}

// GetQuiz は、場所のクイズを取得する
// GET /api/session/{session_id}/s6/quiz/{place_id}
func (h *S6Handler) GetQuiz(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionIDStr := r.PathValue("session_id")
	placeIDStr := r.PathValue("place_id")

	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	placeID, err := valueobject.NewPlaceID(placeIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	quiz, err := h.s6UseCase.GetQuiz(ctx, sessionID, placeID)
	if err != nil {
		handleError(w, err)
		return
	}

	response := dto.ToQuizResponse(quiz)
	respondJSON(w, http.StatusOK, response)
}

// SubmitAnswer は、クイズの回答を送信する
// POST /api/session/{session_id}/s6/answer
func (h *S6Handler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionIDStr := r.PathValue("session_id")

	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	var req struct {
		PlaceID     string `json:"place_id"`
		AnswerIndex int    `json:"answer_index"`
	}
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

	correct, err := h.s6UseCase.SubmitAnswer(ctx, sessionID, placeID, req.AnswerIndex)
	if err != nil {
		handleError(w, err)
		return
	}

	response := dto.QuizAnswerResponse{Correct: correct}
	respondJSON(w, http.StatusOK, response)
}

// GetProgress は、S6の進捗を取得する
// GET /api/session/{session_id}/s6/progress
func (h *S6Handler) GetProgress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionIDStr := r.PathValue("session_id")

	sessionID, err := valueobject.NewSessionIDFromString(sessionIDStr)
	if err != nil {
		handleError(w, err)
		return
	}

	progressList, err := h.s6UseCase.GetProgress(ctx, sessionID)
	if err != nil {
		handleError(w, err)
		return
	}

	responses := dto.ToS6ProgressResponses(progressList)
	respondJSON(w, http.StatusOK, responses)
}
