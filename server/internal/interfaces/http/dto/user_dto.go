package dto

import (
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/model"
)

type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Count int            `json:"count"`
}

func ToUserResponse(user *model.User) UserResponse {
	return UserResponse{
		ID:        user.ID().Value(),
		Name:      user.Name(),
		Email:     user.Email().Value(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

func ToUsersResponse(users []*model.User) UsersResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(user)
	}

	return UsersResponse{
		Users: userResponses,
		Count: len(userResponses),
	}
}
