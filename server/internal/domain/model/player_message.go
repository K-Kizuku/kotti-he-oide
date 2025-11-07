package model

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// PlayerMessage は、プレイヤーが刻んだメッセージを表すドメインモデル
type PlayerMessage struct {
	ID          int
	SessionID   valueobject.SessionID
	PlaceID     valueobject.PlaceID
	MessageText string
	CreatedAt   time.Time
}

const (
	MaxMessageLength = 120 // 最大120文字
)

// NewPlayerMessage は、新しいプレイヤーメッセージを作成する
func NewPlayerMessage(
	sessionID valueobject.SessionID,
	placeID valueobject.PlaceID,
	messageText string,
) (*PlayerMessage, error) {
	// メッセージの長さをチェック
	if len([]rune(messageText)) > MaxMessageLength {
		return nil, errors.NewDomainError(
			errors.INVALID_INPUT,
			"message text exceeds maximum length of 120 characters",
			nil,
		)
	}

	// 空のメッセージは許可しない
	if messageText == "" {
		return nil, errors.NewDomainError(
			errors.INVALID_INPUT,
			"message text cannot be empty",
			nil,
		)
	}

	return &PlayerMessage{
		SessionID:   sessionID,
		PlaceID:     placeID,
		MessageText: messageText,
		CreatedAt:   time.Now(),
	}, nil
}
