package dto

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
)

// MessageRequest は、メッセージ保存のリクエストDTO
type MessageRequest struct {
	PlaceID     string `json:"place_id"`
	MessageText string `json:"message_text"`
}

// MessageResponse は、メッセージ情報のレスポンスDTO
type MessageResponse struct {
	ID          int       `json:"id"`
	PlaceID     string    `json:"place_id"`
	MessageText string    `json:"message_text"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToMessageResponse は、PlayerMessageモデルからMessageResponseDTOに変換する
func ToMessageResponse(message *model.PlayerMessage) *MessageResponse {
	return &MessageResponse{
		ID:          message.ID,
		PlaceID:     message.PlaceID.String(),
		MessageText: message.MessageText,
		CreatedAt:   message.CreatedAt,
	}
}

// ToMessageResponses は、PlayerMessageモデルのスライスからMessageResponseDTOのスライスに変換する
func ToMessageResponses(messages []*model.PlayerMessage) []*MessageResponse {
	responses := make([]*MessageResponse, len(messages))
	for i, message := range messages {
		responses[i] = ToMessageResponse(message)
	}
	return responses
}
