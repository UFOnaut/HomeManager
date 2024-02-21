package errors

import "fmt"

type GeneralError struct {
	Message string
}

func (e *GeneralError) Error() string {
	return fmt.Sprintf("General error: %s", e.Message)
}
