package service

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// SessionService は、セッション管理に関するドメインサービス
type SessionService struct {
	sessionRepo repository.SessionRepository
}

// NewSessionService は、新しいSessionServiceを作成する
func NewSessionService(sessionRepo repository.SessionRepository) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
	}
}

// ValidateSession は、セッションの有効性を検証する
func (s *SessionService) ValidateSession(ctx context.Context, sessionID valueobject.SessionID) (*model.Session, error) {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, errors.NewDomainError(
			errors.SESSION_NOT_FOUND,
			"session not found",
			err,
		)
	}

	if session.IsExpired() {
		return nil, errors.NewDomainError(
			errors.SESSION_EXPIRED,
			"session has expired",
			nil,
		)
	}

	return session, nil
}
