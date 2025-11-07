package model

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// Session は、ゲームセッションを表すドメインモデル
type Session struct {
	SessionID    valueobject.SessionID
	CurrentScene string
	S6StartedAt  *time.Time
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

// NewSession は、新しいセッションを作成する
func NewSession(ttlMinutes int) *Session {
	now := time.Now()
	return &Session{
		SessionID:    valueobject.NewSessionID(),
		CurrentScene: "S0",
		S6StartedAt:  nil,
		CreatedAt:    now,
		ExpiresAt:    now.Add(time.Duration(ttlMinutes) * time.Minute),
	}
}

// IsExpired は、セッションが有効期限切れかどうかを判定する
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// StartS6 は、S6を開始する（7分タイマーのスタート時刻を記録）
func (s *Session) StartS6() {
	now := time.Now()
	s.S6StartedAt = &now
	s.CurrentScene = "S6"
}

// IsS6TimeExpired は、S6の7分制限時間が切れているかどうかを判定する
func (s *Session) IsS6TimeExpired() bool {
	if s.S6StartedAt == nil {
		return false
	}
	elapsed := time.Since(*s.S6StartedAt)
	return elapsed > 7*time.Minute
}

// UpdateScene は、現在のシーンを更新する
func (s *Session) UpdateScene(scene string) {
	s.CurrentScene = scene
}
