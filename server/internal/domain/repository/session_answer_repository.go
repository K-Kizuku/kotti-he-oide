package repository

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// SessionAnswerRepository は、セッション回答の永続化を担当するリポジトリインターフェース
type SessionAnswerRepository interface {
	// Save は、回答を保存する（既存の回答がある場合は上書き）
	Save(ctx context.Context, answer *model.SessionAnswer) error

	// FindBySessionID は、セッションIDで全回答を取得する
	FindBySessionID(ctx context.Context, sessionID valueobject.SessionID) ([]*model.SessionAnswer, error)

	// FindBySessionIDAndQuestionID は、セッションIDと質問IDで回答を取得する
	FindBySessionIDAndQuestionID(
		ctx context.Context,
		sessionID valueobject.SessionID,
		questionID valueobject.QuestionID,
	) (*model.SessionAnswer, error)

	// GetRandomAnswers は、ランダムな過去回答を取得する（クイズのダミー選択肢用）
	GetRandomAnswers(ctx context.Context, questionID valueobject.QuestionID, limit int) ([]string, error)
}
