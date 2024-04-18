package domainerrors

import (
	"errors"
	"strings"
)

// ValidationError is a validation error
type ValidationError struct {
	errors []error
}

func NewValidationError() *ValidationError {
	return &ValidationError{}
}

// IsSuccess returns true if there are no errors
func (e *ValidationError) IsSuccess() bool {
	return len(e.errors) == 0
}

// IsFailure returns true if there are errors
func (e ValidationError) IsFailure() bool {
	return !e.IsSuccess()
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.IsSuccess() {
		return ""
	}

	s := make([]string, len(e.errors))

	for i, err := range e.errors {
		s[i] = err.Error()
	}

	return strings.Join(s, ";")
}

// AddErrorMsg adds an error message to the validation error
func (e *ValidationError) AddErrorMsg(msg string) {
	if strings.TrimSpace(msg) == "" {
		return
	}
	err := errors.New(msg)
	e.AddError(err)
}

// AddError adds an error to the validation error
func (e *ValidationError) AddError(err error) {
	if err == nil {
		return
	}
	e.errors = append(e.errors, err)
}
