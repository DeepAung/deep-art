package handlers

import "github.com/DeepAung/deep-art/api/services"

type UsersHandler struct {
	usersSvc services.UsersSvc
}

func NewUsersHandler(usersSvc services.UsersSvc) *UsersHandler {
	return &UsersHandler{
		usersSvc: usersSvc,
	}
}
