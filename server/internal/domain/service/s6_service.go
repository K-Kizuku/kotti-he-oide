package service

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// S6Service は、S6進捗管理に関するドメインサービス
type S6Service struct {
	s6ProgressRepo repository.S6ProgressRepository
	sessionRepo    repository.SessionRepository
}

// NewS6Service は、新しいS6Serviceを作成する
func NewS6Service(
	s6ProgressRepo repository.S6ProgressRepository,
	sessionRepo repository.SessionRepository,
) *S6Service {
	return &S6Service{
		s6ProgressRepo: s6ProgressRepo,
		sessionRepo:    sessionRepo,
	}
}

// ValidateS6Time は、S6の制限時間（7分）をチェックする
func (s *S6Service) ValidateS6Time(ctx context.Context, sessionID valueobject.SessionID) error {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return errors.NewDomainError(
			errors.SESSION_NOT_FOUND,
			"session not found",
			err,
		)
	}

	if session.S6StartedAt == nil {
		return errors.NewDomainError(
			errors.S6_NOT_STARTED,
			"S6 has not started yet",
			nil,
		)
	}

	if session.IsS6TimeExpired() {
		return errors.NewDomainError(
			errors.S6_TIME_EXPIRED,
			"S6 time limit (7 minutes) has expired",
			nil,
		)
	}

	return nil
}

// IsAllPlacesCompleted は、5箇所すべてが完了したかチェックする
func (s *S6Service) IsAllPlacesCompleted(ctx context.Context, sessionID valueobject.SessionID) (bool, error) {
	count, err := s.s6ProgressRepo.CountCompleted(ctx, sessionID)
	if err != nil {
		return false, err
	}
	return count == 5, nil
}
