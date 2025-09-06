package valueobject

import (
	"fmt"
	"strconv"
)

type JobID struct {
	value int64
}

func NewJobID(value int64) (JobID, error) {
	if value <= 0 {
		return JobID{}, fmt.Errorf("job ID must be positive")
	}
	return JobID{value: value}, nil
}

func JobIDFromString(s string) (JobID, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return JobID{}, fmt.Errorf("invalid job ID format: %w", err)
	}
	return NewJobID(value)
}

func (id JobID) Value() int64 {
	return id.value
}

func (id JobID) String() string {
	return strconv.FormatInt(id.value, 10)
}

func (id JobID) Equals(other JobID) bool {
	return id.value == other.value
}