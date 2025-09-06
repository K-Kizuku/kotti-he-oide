package valueobject

import (
	"fmt"
	"strconv"
)

type SubscriptionID struct {
	value int64
}

func NewSubscriptionID(value int64) (SubscriptionID, error) {
	if value <= 0 {
		return SubscriptionID{}, fmt.Errorf("subscription ID must be positive")
	}
	return SubscriptionID{value: value}, nil
}

func SubscriptionIDFromString(s string) (SubscriptionID, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return SubscriptionID{}, fmt.Errorf("invalid subscription ID format: %w", err)
	}
	return NewSubscriptionID(value)
}

func (id SubscriptionID) Value() int64 {
	return id.value
}

func (id SubscriptionID) String() string {
	return strconv.FormatInt(id.value, 10)
}

func (id SubscriptionID) Equals(other SubscriptionID) bool {
	return id.value == other.value
}
