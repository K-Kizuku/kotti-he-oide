package valueobject

import (
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
	"github.com/google/uuid"
)

// SessionID は、セッションIDを表すValue Object
type SessionID struct {
	value uuid.UUID
}

// NewSessionID は、新しいSessionIDを生成する
func NewSessionID() SessionID {
	return SessionID{value: uuid.New()}
}

// NewSessionIDFromString は、文字列からSessionIDを作成する
func NewSessionIDFromString(id string) (SessionID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return SessionID{}, errors.NewDomainError(
			errors.INVALID_INPUT,
			"invalid session ID format",
			err,
		)
	}
	return SessionID{value: parsed}, nil
}

// String は、SessionIDを文字列として返す
func (s SessionID) String() string {
	return s.value.String()
}

// UUID は、SessionIDをUUIDとして返す
func (s SessionID) UUID() uuid.UUID {
	return s.value
}

// Equals は、2つのSessionIDが等しいかどうかを判定する
func (s SessionID) Equals(other SessionID) bool {
	return s.value == other.value
}
