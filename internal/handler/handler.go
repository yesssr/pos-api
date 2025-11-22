package handler

import "pos-api/internal/store"

type Handler struct {
	Auth *AuthHandler
	User *UserHandler
}

func New(q *store.Queries) *Handler {
	return &Handler{
		Auth: NewAuthHandler(q),
		User: NewUserHandler(q),
	}
}
