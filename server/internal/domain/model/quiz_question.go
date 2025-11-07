package model

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// QuizQuestion は、クイズ問題を表すドメインモデル
type QuizQuestion struct {
	QuizID       valueobject.QuizID
	SessionID    valueobject.SessionID
	PlaceID      valueobject.PlaceID
	QuestionText string
	Options      [4]string // 4択の選択肢
	AnswerIndex  int       // 正解のインデックス（0-3）
	CreatedAt    time.Time
}

// NewQuizQuestion は、新しいクイズ問題を作成する
func NewQuizQuestion(
	sessionID valueobject.SessionID,
	placeID valueobject.PlaceID,
	questionText string,
	options [4]string,
	answerIndex int,
) *QuizQuestion {
	return &QuizQuestion{
		QuizID:       valueobject.NewQuizID(),
		SessionID:    sessionID,
		PlaceID:      placeID,
		QuestionText: questionText,
		Options:      options,
		AnswerIndex:  answerIndex,
		CreatedAt:    time.Now(),
	}
}

// CheckAnswer は、回答が正解かどうかをチェックする
func (q *QuizQuestion) CheckAnswer(answerIndex int) bool {
	return answerIndex == q.AnswerIndex
}

// GetCorrectAnswer は、正解の選択肢を返す
func (q *QuizQuestion) GetCorrectAnswer() string {
	return q.Options[q.AnswerIndex]
}
