package dto

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
)

// SessionResponse は、セッション情報のレスポンスDTO
type SessionResponse struct {
	SessionID    string     `json:"session_id"`
	CurrentScene string     `json:"current_scene"`
	S6StartedAt  *time.Time `json:"s6_started_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	ExpiresAt    time.Time  `json:"expires_at"`
}

// ToSessionResponse は、SessionモデルからSessionResponseDTOに変換する
func ToSessionResponse(session *model.Session) *SessionResponse {
	return &SessionResponse{
		SessionID:    session.SessionID.String(),
		CurrentScene: session.CurrentScene,
		S6StartedAt:  session.S6StartedAt,
		CreatedAt:    session.CreatedAt,
		ExpiresAt:    session.ExpiresAt,
	}
}

// UpdateSceneRequest は、シーン更新のリクエストDTO
type UpdateSceneRequest struct {
	Scene string `json:"scene"`
}
