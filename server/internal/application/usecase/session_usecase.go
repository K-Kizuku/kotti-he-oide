package usecase

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

// SessionUseCase は、セッション管理のユースケース
type SessionUseCase struct {
	sessionRepo    repository.SessionRepository
	sessionService *service.SessionService
	sessionTTL     int // セッションの有効期限（分）
}

// NewSessionUseCase は、新しいSessionUseCaseを作成する
func NewSessionUseCase(
	sessionRepo repository.SessionRepository,
	sessionService *service.SessionService,
	sessionTTL int,
) *SessionUseCase {
	return &SessionUseCase{
		sessionRepo:    sessionRepo,
		sessionService: sessionService,
		sessionTTL:     sessionTTL,
	}
}

// CreateSession は、新しいセッションを作成する
func (u *SessionUseCase) CreateSession(ctx context.Context) (*model.Session, error) {
	session := model.NewSession(u.sessionTTL)

	if err := u.sessionRepo.Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession は、セッションIDでセッションを取得する
func (u *SessionUseCase) GetSession(ctx context.Context, sessionID valueobject.SessionID) (*model.Session, error) {
	return u.sessionService.ValidateSession(ctx, sessionID)
}

// UpdateScene は、現在のシーンを更新する
func (u *SessionUseCase) UpdateScene(ctx context.Context, sessionID valueobject.SessionID, scene string) error {
	session, err := u.sessionService.ValidateSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.UpdateScene(scene)

	return u.sessionRepo.Update(ctx, session)
}

// StartS6 は、S6を開始する（7分タイマー開始）
func (u *SessionUseCase) StartS6(ctx context.Context, sessionID valueobject.SessionID) error {
	session, err := u.sessionService.ValidateSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.StartS6()

	return u.sessionRepo.Update(ctx, session)
}
