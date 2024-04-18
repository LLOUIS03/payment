package domainerrors

// NoTransactionError is an error type for invalid credentials
type NoTransactionError struct {
	base
}

// Message returns the error message
func (e NoTransactionError) Message() string {
	return e.message
}

// NewNoTransactionError creates a new NoTransactionError
func NewNoTransactionError(innerError error) NoTransactionError {
	return NoTransactionError{
		base: base{
			message:    "no transaction found",
			innerError: innerError,
		},
	}
}
