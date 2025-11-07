package model

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// S6Progress は、S6の進捗状況を表すドメインモデル
type S6Progress struct {
	ID         int
	SessionID  valueobject.SessionID
	PlaceID    valueobject.PlaceID
	Verified   bool
	VerifiedBy string // "photo" or "manual"
	QuizID     *valueobject.QuizID
	Answered   bool
	Correct    bool
	VerifiedAt *time.Time
}

// NewS6Progress は、新しいS6進捗を作成する
func NewS6Progress(
	sessionID valueobject.SessionID,
	placeID valueobject.PlaceID,
) *S6Progress {
	return &S6Progress{
		SessionID:  sessionID,
		PlaceID:    placeID,
		Verified:   false,
		VerifiedBy: "",
		QuizID:     nil,
		Answered:   false,
		Correct:    false,
		VerifiedAt: nil,
	}
}

// MarkVerified は、場所到達を記録する
func (p *S6Progress) MarkVerified(verifiedBy string) {
	now := time.Now()
	p.Verified = true
	p.VerifiedBy = verifiedBy
	p.VerifiedAt = &now
}

// SetQuiz は、クイズIDを設定する
func (p *S6Progress) SetQuiz(quizID valueobject.QuizID) {
	p.QuizID = &quizID
}

// MarkAnswered は、クイズに回答したことを記録する
func (p *S6Progress) MarkAnswered(correct bool) {
	p.Answered = true
	p.Correct = correct
}

// IsCompleted は、この場所が完了したかどうかを判定する
func (p *S6Progress) IsCompleted() bool {
	return p.Verified && p.Answered && p.Correct
}
