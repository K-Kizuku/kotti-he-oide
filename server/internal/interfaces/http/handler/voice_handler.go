package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/K-Kizuku/kotti-he-oide/internal/application/usecase"
	"github.com/K-Kizuku/kotti-he-oide/internal/interfaces/http/dto"
)

// VoiceHandler は、音声生成に関するHTTPハンドラー
type VoiceHandler struct {
	voiceUseCase *usecase.VoiceUseCase
}

// NewVoiceHandler は、新しいVoiceHandlerを作成する
func NewVoiceHandler(voiceUseCase *usecase.VoiceUseCase) *VoiceHandler {
	return &VoiceHandler{
		voiceUseCase: voiceUseCase,
	}
}

// GenerateVoice は、音声生成リクエストを処理する
// POST /api/voice/generate
func (h *VoiceHandler) GenerateVoice(w http.ResponseWriter, r *http.Request) {
	// リクエストボディをデコード
	var req dto.GenerateVoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[VoiceHandler] Failed to decode request: %v", err)
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// テキストのバリデーション
	if strings.TrimSpace(req.Text) == "" {
		h.respondError(w, http.StatusBadRequest, "EMPTY_TEXT", "Text is required", "")
		return
	}

	// 音声を生成
	audioURL, err := h.voiceUseCase.GenerateVoice(req.Text, req.SpeakerID)
	if err != nil {
		log.Printf("[VoiceHandler] Failed to generate voice: %v", err)
		h.respondError(w, http.StatusInternalServerError, "VOICE_GENERATION_FAILED", "Failed to generate voice", err.Error())
		return
	}

	// レスポンスを返す
	resp := dto.GenerateVoiceResponse{
		AudioURL: audioURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("[VoiceHandler] Failed to encode response: %v", err)
	}
}

// respondError は、エラーレスポンスを返す
func (h *VoiceHandler) respondError(w http.ResponseWriter, statusCode int, errorCode, message, details string) {
	resp := dto.ErrorResponse{
		Error:   errorCode,
		Message: message,
		Details: details,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("[VoiceHandler] Failed to encode error response: %v", err)
	}
}
