package domainerrors

// UnauthorizedTranError is an error type for unauthorized transactions
type UnauthorizedTranError struct {
	base
}

// Message returns the error message
func (e UnauthorizedTranError) Message() string {
	return e.message
}

// NewUnauthorizedTranError creates a new UnauthorizedTranError
func NewUnauthorizedTranError(innerError error) UnauthorizedTranError {
	return UnauthorizedTranError{
		base: base{
			message:    "unable to authorize transaction",
			innerError: innerError,
		},
	}
}
