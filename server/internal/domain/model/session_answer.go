package model

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// SessionAnswer は、S4の質問への回答を表すドメインモデル
type SessionAnswer struct {
	ID         int
	SessionID  valueobject.SessionID
	QuestionID valueobject.QuestionID
	AnswerText string
	AnsweredAt time.Time
}

// NewSessionAnswer は、新しい回答を作成する
func NewSessionAnswer(
	sessionID valueobject.SessionID,
	questionID valueobject.QuestionID,
	answerText string,
) *SessionAnswer {
	return &SessionAnswer{
		SessionID:  sessionID,
		QuestionID: questionID,
		AnswerText: answerText,
		AnsweredAt: time.Now(),
	}
}
