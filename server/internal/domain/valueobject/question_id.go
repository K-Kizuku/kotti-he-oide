package valueobject

import (
	"github.com/K-Kizuku/kotti-he-oide/pkg/errors"
)

// QuestionID は、S4の質問IDを表すValue Object（1-10）
type QuestionID struct {
	value int
}

// NewQuestionID は、整数からQuestionIDを作成する
func NewQuestionID(id int) (QuestionID, error) {
	if id < 1 || id > 10 {
		return QuestionID{}, errors.NewDomainError(
			errors.INVALID_INPUT,
			"question ID must be between 1 and 10",
			nil,
		)
	}
	return QuestionID{value: id}, nil
}

// Int は、QuestionIDを整数として返す
func (q QuestionID) Int() int {
	return q.value
}

// Equals は、2つのQuestionIDが等しいかどうかを判定する
func (q QuestionID) Equals(other QuestionID) bool {
	return q.value == other.value
}
