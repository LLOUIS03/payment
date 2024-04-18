package domainerrors

// CredentialsExpiredError is an error type for expired credentials
type CredentialsExpiredError struct {
	base
}

// Message returns the error message
func (e CredentialsExpiredError) Message() string {
	return e.message
}

// NewCredentialsExpiredError creates a new CredentialsExpiredError
func NewCredentialsExpiredError(innerError error) CredentialsExpiredError {
	return CredentialsExpiredError{
		base: base{
			message:    "credentials expired, please login again",
			innerError: innerError,
		},
	}
}
