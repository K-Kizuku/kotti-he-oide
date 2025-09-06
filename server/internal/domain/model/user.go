package model

import (
	"fmt"
	"time"

	"github.com/K-Kizuku/kotti-he-oide/internal/domain/valueobject"
)

type User struct {
	id        valueobject.UserID
	name      string
	email     valueobject.Email
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(id valueobject.UserID, name string, email valueobject.Email) *User {
	now := time.Now()
	return &User{
		id:        id,
		name:      name,
		email:     email,
		createdAt: now,
		updatedAt: now,
	}
}

func ReconstructUser(id valueobject.UserID, name string, email valueobject.Email, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		name:      name,
		email:     email,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) ID() valueobject.UserID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) ChangeName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	u.name = name
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ChangeEmail(email valueobject.Email) {
	u.email = email
	u.updatedAt = time.Now()
}