package usecase

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// S6UseCase は、S6進捗管理のユースケース
type S6UseCase struct {
	s6ProgressRepo repository.S6ProgressRepository
	quizService    *service.QuizService
	s6Service      *service.S6Service
	sessionService *service.SessionService
}

// NewS6UseCase は、新しいS6UseCaseを作成する
func NewS6UseCase(
	s6ProgressRepo repository.S6ProgressRepository,
	quizService *service.QuizService,
	s6Service *service.S6Service,
	sessionService *service.SessionService,
) *S6UseCase {
	return &S6UseCase{
		s6ProgressRepo: s6ProgressRepo,
		quizService:    quizService,
		s6Service:      s6Service,
		sessionService: sessionService,
	}
}

// InitializeProgress は、5箇所の進捗を初期化する
func (u *S6UseCase) InitializeProgress(ctx context.Context, sessionID valueobject.SessionID) error {
	// セッションの有効性をチェック
	if _, err := u.sessionService.ValidateSession(ctx, sessionID); err != nil {
		return err
	}

	// 5箇所分の進捗を作成
	allPlaces := valueobject.GetAllPlaces()
	for _, placeStr := range allPlaces {
		placeID, err := valueobject.NewPlaceID(placeStr)
		if err != nil {
			return err
		}

		progress := model.NewS6Progress(sessionID, placeID)
		if err := u.s6ProgressRepo.Save(ctx, progress); err != nil {
			return err
		}
	}

	return nil
}

// MarkPlaceVerified は、場所到達を記録する
func (u *S6UseCase) MarkPlaceVerified(
	ctx context.Context,
	sessionID valueobject.SessionID,
	placeID valueobject.PlaceID,
	verifiedBy string,
) error {
	// セッションの有効性をチェック
	if _, err := u.sessionService.ValidateSession(ctx, sessionID); err != nil {
		return err
	}

	// S6の時間制限をチェック
	if err := u.s6Service.ValidateS6Time(ctx, sessionID); err != nil {
		return err
	}

	// 進捗を取得
	progress, err := u.s6ProgressRepo.FindBySessionIDAndPlaceID(ctx, sessionID, placeID)
	if err != nil {
		return err
	}

	// すでに完了している場合はエラー
	if progress.IsCompleted() {
		return errors.NewDomainError(
			errors.S6_ALREADY_COMPLETED,
			"this place is already completed",
			nil,
		)
	}

	// 到達を記録
	progress.MarkVerified(verifiedBy)

	return u.s6ProgressRepo.Update(ctx, progress)
}

// GetQuiz は、場所のクイズを取得する（なければ生成）
func (u *S6UseCase) GetQuiz(
	ctx context.Context,
	sessionID valueobject.SessionID,
	placeID valueobject.PlaceID,
) (*model.QuizQuestion, error) {
	// セッションの有効性をチェック
	if _, err := u.sessionService.ValidateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	// S6の時間制限をチェック
	if err := u.s6Service.ValidateS6Time(ctx, sessionID); err != nil {
		return nil, err
	}

	// 進捗を確認（場所に到達しているか）
	progress, err := u.s6ProgressRepo.FindBySessionIDAndPlaceID(ctx, sessionID, placeID)
	if err != nil {
		return nil, err
	}

	if !progress.Verified {
		return nil, errors.NewDomainError(
			errors.INVALID_INPUT,
			"place not verified yet",
			nil,
		)
	}

	// クイズを生成（または既存のものを取得）
	quiz, err := u.quizService.GenerateQuiz(ctx, sessionID, placeID)
	if err != nil {
		return nil, err
	}

	// 進捗にクイズIDを設定
	if progress.QuizID == nil {
		progress.SetQuiz(quiz.QuizID)
		if err := u.s6ProgressRepo.Update(ctx, progress); err != nil {
			return nil, err
		}
	}

	return quiz, nil
}

// SubmitAnswer は、クイズ回答を送信する
func (u *S6UseCase) SubmitAnswer(
	ctx context.Context,
	sessionID valueobject.SessionID,
	placeID valueobject.PlaceID,
	answerIndex int,
) (bool, error) {
	// セッションの有効性をチェック
	if _, err := u.sessionService.ValidateSession(ctx, sessionID); err != nil {
		return false, err
	}

	// S6の時間制限をチェック（最後のピース回答中は猶予）
	progress, err := u.s6ProgressRepo.FindBySessionIDAndPlaceID(ctx, sessionID, placeID)
	if err != nil {
		return false, err
	}

	// 進捗を確認
	if !progress.Verified {
		return false, errors.NewDomainError(
			errors.INVALID_INPUT,
			"place not verified yet",
			nil,
		)
	}

	if progress.QuizID == nil {
		return false, errors.NewDomainError(
			errors.QUIZ_NOT_FOUND,
			"quiz not generated yet",
			nil,
		)
	}

	// クイズを取得
	quiz, err := u.quizService.GenerateQuiz(ctx, sessionID, placeID)
	if err != nil {
		return false, err
	}

	// 回答をチェック
	correct := quiz.CheckAnswer(answerIndex)

	// 進捗を更新
	progress.MarkAnswered(correct)
	if err := u.s6ProgressRepo.Update(ctx, progress); err != nil {
		return false, err
	}

	return correct, nil
}

// GetProgress は、セッションの全進捗を取得する
func (u *S6UseCase) GetProgress(ctx context.Context, sessionID valueobject.SessionID) ([]*model.S6Progress, error) {
	// セッションの有効性をチェック
	if _, err := u.sessionService.ValidateSession(ctx, sessionID); err != nil {
		return nil, err
	}

	return u.s6ProgressRepo.FindBySessionID(ctx, sessionID)
}

// IsAllCompleted は、すべての場所が完了したかチェックする
func (u *S6UseCase) IsAllCompleted(ctx context.Context, sessionID valueobject.SessionID) (bool, error) {
	// セッションの有効性をチェック
	if _, err := u.sessionService.ValidateSession(ctx, sessionID); err != nil {
		return false, err
	}

	return u.s6Service.IsAllPlacesCompleted(ctx, sessionID)
}
