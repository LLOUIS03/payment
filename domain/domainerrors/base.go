package domainerrors

import "fmt"

type base struct {
	message    string
	innerError error
}

// Message returns the error message
func (e base) Message() string {
	return e.message
}

// Error returns the error message
func (e base) Error() string {
	if e.innerError != nil {
		return fmt.Sprintf("%s:%s", e.message, e.innerError.Error())
	}

	return e.message
}
