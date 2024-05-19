package httperror

import (
	"fmt"
)

type HttpError struct {
	Msg    string
	Status int
}

func New(msg string, status int) error {
	return &HttpError{
		Msg:    msg,
		Status: status,
	}
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Msg)
}
