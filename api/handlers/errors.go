package handlers

import "errors"

var (
	ErrPayloadNotFound  = errors.New("payload from middleware not found")
	ErrUserDataNotFound = errors.New("user data from middleware not found")
)
