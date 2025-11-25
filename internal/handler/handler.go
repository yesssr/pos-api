package handler

import (
	"pos-api/internal/store"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Handler struct {
	R2Client *s3.Client;
	Auth *AuthHandler;
	User *UserHandler;
}

func New(q *store.Queries, r *s3.Client) *Handler {
	return &Handler{
		Auth: NewAuthHandler(q),
		User: NewUserHandler(q, r),
	}
}
