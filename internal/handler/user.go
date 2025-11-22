package handler

import (
	"net/http"
	"pos-api/internal/lib"
	"pos-api/internal/store"
)

type CreateUserInput struct {
  Username string `json:"username" validate:"required,min=3"`
  Password string `json:"password" validate:"required,min=6"`
  Role     store.Roles `json:"role" validate:"required,oneof=admin kasir"`
}

type UserHandler struct {
	queries *store.Queries
}

func NewUserHandler(q *store.Queries) *UserHandler {
	return &UserHandler{
		queries: q,
	}
}

func (s *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var b CreateUserInput;
	if !lib.ValidateJSON(w, r, &b) {
		return
	}

	args := store.CreateUserParams{
		Username: b.Username,
		Password: b.Password,
		Role:     b.Role,
		ImageUrl: "",
	}

	u, err := s.queries.CreateUser(r.Context(), args);
	if err != nil {
		lib.SendErrorResponse(w, err)
		return
	}

	lib.SendResponse(w, true, http.StatusCreated, "Successfully added user", u, nil, nil)
}
