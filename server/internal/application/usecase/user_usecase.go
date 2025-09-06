package usecase

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/repository"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/service"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

type UserUseCase struct {
	userRepo    repository.UserRepository
	userService *service.UserService
}

func NewUserUseCase(userRepo repository.UserRepository, userService *service.UserService) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		userService: userService,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, name, emailStr string) (*model.User, error) {
	email, err := valueobject.NewEmail(emailStr)
	if err != nil {
		return nil, errors.WrapDomainError(errors.ErrInvalidEmail.Code, "Invalid email", err)
	}

	if err := u.userService.CheckEmailDuplicate(ctx, email); err != nil {
		return nil, errors.WrapDomainError(errors.ErrEmailAlreadyExist.Code, "Email already exists", err)
	}

	userID, err := u.userRepo.NextIdentity(ctx)
	if err != nil {
		return nil, err
	}

	user := model.NewUser(userID, name, email)

	if err := u.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) GetUser(ctx context.Context, userIDInt int) (*model.User, error) {
	userID, err := valueobject.NewUserID(userIDInt)
	if err != nil {
		return nil, errors.WrapDomainError(errors.ErrInvalidUserID.Code, "Invalid user ID", err)
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func (u *UserUseCase) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	return u.userRepo.FindAll(ctx)
}

func (u *UserUseCase) DeleteUser(ctx context.Context, userIDInt int) error {
	userID, err := valueobject.NewUserID(userIDInt)
	if err != nil {
		return errors.WrapDomainError(errors.ErrInvalidUserID.Code, "Invalid user ID", err)
	}

	if err := u.userService.CanDeleteUser(ctx, userID); err != nil {
		return err
	}

	return u.userRepo.Delete(ctx, userID)
}