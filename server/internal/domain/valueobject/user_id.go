package valueobject

import (
	"fmt"
	"strconv"
)

type UserID struct {
	value int
}

func NewUserID(value int) (UserID, error) {
	if value <= 0 {
		return UserID{}, fmt.Errorf("user ID must be positive, got %d", value)
	}
	return UserID{value: value}, nil
}

func UserIDFromString(s string) (UserID, error) {
	value, err := strconv.Atoi(s)
	if err != nil {
		return UserID{}, fmt.Errorf("invalid user ID format: %w", err)
	}
	return NewUserID(value)
}

func (u UserID) Value() int {
	return u.value
}

func (u UserID) String() string {
	return strconv.Itoa(u.value)
}

func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}
