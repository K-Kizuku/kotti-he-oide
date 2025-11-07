package valueobject

import (
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
	"github.com/google/uuid"
)

// SubscriptionID は、プッシュ通知サブスクリプションIDを表すValue Object
type SubscriptionID struct {
	value uuid.UUID
}

// NewSubscriptionID は、新しいSubscriptionIDを生成する
func NewSubscriptionID() SubscriptionID {
	return SubscriptionID{value: uuid.New()}
}

// NewSubscriptionIDFromString は、文字列からSubscriptionIDを作成する
func NewSubscriptionIDFromString(id string) (SubscriptionID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return SubscriptionID{}, errors.NewDomainError(
			errors.INVALID_INPUT,
			"invalid subscription ID format",
			err,
		)
	}
	return SubscriptionID{value: parsed}, nil
}

// String は、SubscriptionIDを文字列として返す
func (s SubscriptionID) String() string {
	return s.value.String()
}

// UUID は、SubscriptionIDをUUIDとして返す
func (s SubscriptionID) UUID() uuid.UUID {
	return s.value
}

// Equals は、2つのSubscriptionIDが等しいかどうかを判定する
func (s SubscriptionID) Equals(other SubscriptionID) bool {
	return s.value == other.value
}
