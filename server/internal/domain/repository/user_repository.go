package repository

import (
	"context"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type UserRepository interface {
	Save(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id valueobject.UserID) (*model.User, error)
	FindByEmail(ctx context.Context, email valueobject.Email) (*model.User, error)
	FindAll(ctx context.Context) ([]*model.User, error)
	Delete(ctx context.Context, id valueobject.UserID) error
	NextIdentity(ctx context.Context) (valueobject.UserID, error)
}