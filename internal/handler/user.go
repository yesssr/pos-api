package handler

import (
	"net/http"
	"pos-api/internal/lib"
	"pos-api/internal/store"
)

type CreateUserInput struct {
  Username string `json:"username" validate:"required,min=3"`
  Password string `json:"password" validate:"required,min=6"`
  Role     store.Roles `json:"role" validate:"required,oneof=admin cashier"`
}

type UserHandler struct {
	queries *store.Queries
}

func NewUserHandler(q *store.Queries) *UserHandler {
	return &UserHandler{
		queries: q,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var b CreateUserInput;
	if !lib.ValidateJSON(w, r, &b) {
		return
	}

	if err := lib.ValidateStruct(&b); err != nil {
		lib.SendErrorResponse(w, err)
		return
	}

	ctx := r.Context()
	if _, err := h.queries.GetUserByUsername(ctx, b.Username); err == nil {
		lib.SendErrorResponse(w, &lib.AppError{
			Message: "Username already exists",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	pass, _ := lib.HashPassword(b.Password);

	args := store.CreateUserParams{
		Username: b.Username,
		Password: pass,
		Role:     b.Role,
		ImageUrl: "",
	}

	u, err := h.queries.CreateUser(ctx, args);
	if err != nil {
		lib.SendErrorResponse(w, err)
		return
	}

	lib.SendResponse(w, http.StatusCreated, "Successfully added user", u, nil, nil)
}
