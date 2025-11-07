package errors

import "fmt"

type DomainError struct {
	Code    string
	Message string
	Cause   error
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewDomainError(code, message string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func WrapDomainError(code, message string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// エラーコード定数
const (
	// 共通エラー
	INVALID_INPUT  = "INVALID_INPUT"
	INVALID_STATE  = "INVALID_STATE"
	NOT_FOUND      = "NOT_FOUND"
	ALREADY_EXISTS = "ALREADY_EXISTS"

	// セッション関連
	SESSION_NOT_FOUND  = "SESSION_NOT_FOUND"
	SESSION_EXPIRED    = "SESSION_EXPIRED"
	INVALID_SESSION_ID = "INVALID_SESSION_ID"

	// 回答関連
	ANSWER_REQUIRED  = "ANSWER_REQUIRED"
	ANSWER_NOT_FOUND = "ANSWER_NOT_FOUND"

	// S6関連
	S6_NOT_STARTED       = "S6_NOT_STARTED"
	S6_TIME_EXPIRED      = "S6_TIME_EXPIRED"
	S6_ALREADY_COMPLETED = "S6_ALREADY_COMPLETED"

	// クイズ関連
	QUIZ_NOT_FOUND = "QUIZ_NOT_FOUND"
	QUIZ_INCORRECT = "QUIZ_INCORRECT"

	// メッセージ関連
	MESSAGE_TOO_LONG = "MESSAGE_TOO_LONG"
	MESSAGE_EMPTY    = "MESSAGE_EMPTY"
)
