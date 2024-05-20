package httperror

import (
	"fmt"
	"net/http"
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

func Extract(err error) (msg string, status int) {
	switch err := err.(type) {
	case *HttpError:
		return err.Msg, err.Status
	default:
		status = http.StatusInternalServerError
		msg = http.StatusText(status)
		return msg, status
	}
}
