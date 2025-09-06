package service

import (
	"context"
	"fmt"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CheckEmailDuplicate(ctx context.Context, email valueobject.Email) error {
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return fmt.Errorf("email already exists: %s", email.Value())
	}
	return nil
}

func (s *UserService) CanDeleteUser(ctx context.Context, userID valueobject.UserID) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found: %s", userID.String())
	}
	return nil
}
