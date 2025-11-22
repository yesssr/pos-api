package handler

import (
	"net/http"
	"pos-api/internal/lib"
	"pos-api/internal/store"
)
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHandler struct {
	queries *store.Queries
}

func NewAuthHandler(q *store.Queries) *AuthHandler {
	return &AuthHandler{ queries: q }
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user LoginRequest;
	if !lib.ValidateJSON(w, r, &user) {
		return
	};
	if user.Username == "" || user.Password == "" {
		lib.SendErrorResponse(w, &lib.AppError{
			Message: "Username and password are required",
			StatusCode: http.StatusBadRequest,
		});
		return;
	}

	u, err := h.queries.GetUserByUsername(r.Context(), user.Username);
	if err != nil {
		lib.SendErrorResponse(w, &lib.AppError{
			Message: err.Error(),
			StatusCode: http.StatusInternalServerError,
		});
		return;
	}

	if err := lib.VerifyPassword(u.Password, user.Password); err != nil {
		lib.SendErrorResponse(w, &lib.AppError{
			Message: "Invalid username or password",
			StatusCode: http.StatusUnauthorized,
		});
		return;
	}

	p := lib.Payload{
		Id: u.ID.String(),
		Username: u.Username,
		Role: string(u.Role),
	}
	c, err := lib.CreateToken(p);
	if err != nil {
		lib.SendErrorResponse(w, err);
		return;
	}
	lib.SendResponse(w, true, http.StatusOK, "Login successful", nil, nil, &c);
}
