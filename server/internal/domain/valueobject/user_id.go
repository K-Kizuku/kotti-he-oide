package valueobject

import "fmt"

type UserID struct {
	value int
}

func NewUserID(value int) (UserID, error) {
	if value <= 0 {
		return UserID{}, fmt.Errorf("user ID must be positive, got %d", value)
	}
	return UserID{value: value}, nil
}

func (u UserID) Value() int {
	return u.value
}

func (u UserID) String() string {
	return fmt.Sprintf("UserID(%d)", u.value)
}

func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}