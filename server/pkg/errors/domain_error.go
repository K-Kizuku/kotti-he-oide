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

func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

func WrapDomainError(code, message string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

var (
	ErrUserNotFound      = NewDomainError("USER_NOT_FOUND", "User not found")
	ErrEmailAlreadyExist = NewDomainError("EMAIL_ALREADY_EXISTS", "Email already exists")
	ErrInvalidUserID     = NewDomainError("INVALID_USER_ID", "Invalid user ID")
	ErrInvalidEmail      = NewDomainError("INVALID_EMAIL", "Invalid email format")
)