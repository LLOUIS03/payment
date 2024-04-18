package domainerrors

// InvalidCredentialsError is an error type for invalid credentials
type InvalidCredentialsError struct {
	base
}

// Message returns the error message
func (e InvalidCredentialsError) Message() string {
	return e.message
}

// NewInvalidCredentialsError creates a new InvalidCredentialsError
func NewInvalidCredentialsError(innerError error) InvalidCredentialsError {
	return InvalidCredentialsError{
		base: base{
			message:    "invalid credentials",
			innerError: innerError,
		},
	}
}
