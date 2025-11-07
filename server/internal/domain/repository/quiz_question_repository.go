package repository

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// QuizQuestionRepository は、クイズ問題の永続化を担当するリポジトリインターフェース
type QuizQuestionRepository interface {
	// Save は、クイズ問題を保存する
	Save(ctx context.Context, quiz *model.QuizQuestion) error

	// FindByID は、クイズIDでクイズ問題を取得する
	FindByID(ctx context.Context, quizID valueobject.QuizID) (*model.QuizQuestion, error)

	// FindBySessionIDAndPlaceID は、セッションIDと場所IDでクイズ問題を取得する
	FindBySessionIDAndPlaceID(
		ctx context.Context,
		sessionID valueobject.SessionID,
		placeID valueobject.PlaceID,
	) (*model.QuizQuestion, error)

	// FindBySessionID は、セッションIDで全クイズ問題を取得する
	FindBySessionID(ctx context.Context, sessionID valueobject.SessionID) ([]*model.QuizQuestion, error)
}
