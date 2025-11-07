package usecase

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// S4AnswerUseCase は、S4回答管理のユースケース
type S4AnswerUseCase struct {
	sessionAnswerRepo repository.SessionAnswerRepository
	sessionService    *service.SessionService
}

// NewS4AnswerUseCase は、新しいS4AnswerUseCaseを作成する
func NewS4AnswerUseCase(
	sessionAnswerRepo repository.SessionAnswerRepository,
	sessionService *service.SessionService,
) *S4AnswerUseCase {
	return &S4AnswerUseCase{
		sessionAnswerRepo: sessionAnswerRepo,
		sessionService:    sessionService,
	}
}

// SaveAnswer は、回答を保存する（逐次保存）
func (u *S4AnswerUseCase) SaveAnswer(
	ctx context.Context,
	sessionID valueobject.SessionID,
	questionID valueobject.QuestionID,
	answerText string,
) error {
	// セッションの有効性をチェック
	if _, err := u.sessionService.ValidateSession(ctx, sessionID); err != nil {
		return err
	}

	// 回答を作成
	answer := model.NewSessionAnswer(sessionID, questionID, answerText)

	// 回答を保存（既存の回答がある場合は上書き）
	return u.sessionAnswerRepo.Save(ctx, answer)
}

// GetAnswers は、セッションIDで全回答を取得する
func (u *S4AnswerUseCase) GetAnswers(ctx context.Context, sessionID valueobject.SessionID) ([]*model.SessionAnswer, error) {
	// セッションの有効性をチェック
	if _, err := u.sessionService.ValidateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	return u.sessionAnswerRepo.FindBySessionID(ctx, sessionID)
}

// GetAnswer は、セッションIDと質問IDで回答を取得する
func (u *S4AnswerUseCase) GetAnswer(
	ctx context.Context,
	sessionID valueobject.SessionID,
	questionID valueobject.QuestionID,
) (*model.SessionAnswer, error) {
	// セッションの有効性をチェック
	if _, err := u.sessionService.ValidateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	return u.sessionAnswerRepo.FindBySessionIDAndQuestionID(ctx, sessionID, questionID)
}
