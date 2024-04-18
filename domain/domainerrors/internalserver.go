package domainerrors

// InternalServerError is an error type for internal server errors
type InternalServerError struct {
	base
}

// Message returns the error message
func (e InternalServerError) Message() string {
	return e.message
}

// NewInternalServerError creates a new InternalServerError
func NewInternalServerError(innerError error) InternalServerError {
	return InternalServerError{
		base: base{
			message:    "internal server error",
			innerError: innerError,
		},
	}
}
