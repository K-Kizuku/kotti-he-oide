package valueobject

import (
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
	"github.com/google/uuid"
)

// QuizID は、クイズIDを表すValue Object
type QuizID struct {
	value uuid.UUID
}

// NewQuizID は、新しいQuizIDを生成する
func NewQuizID() QuizID {
	return QuizID{value: uuid.New()}
}

// NewQuizIDFromString は、文字列からQuizIDを作成する
func NewQuizIDFromString(id string) (QuizID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return QuizID{}, errors.NewDomainError(
			errors.INVALID_INPUT,
			"invalid quiz ID format",
			err,
		)
	}
	return QuizID{value: parsed}, nil
}

// String は、QuizIDを文字列として返す
func (q QuizID) String() string {
	return q.value.String()
}

// UUID は、QuizIDをUUIDとして返す
func (q QuizID) UUID() uuid.UUID {
	return q.value
}

// Equals は、2つのQuizIDが等しいかどうかを判定する
func (q QuizID) Equals(other QuizID) bool {
	return q.value == other.value
}
