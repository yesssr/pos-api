package handler

import (
	"net/http"
	"pos-api/internal/lib"
	"pos-api/internal/service"
)
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHandler struct {
	s *service.AuthService;
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{ s: s }
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user LoginRequest;
	if !lib.ValidateJSON(w, r, &user) {
		return;
	};

	u, c, err := h.s.Login(r.Context(), user.Username, user.Password);
	if err != nil {
		lib.SendErrorResponse(w, err, nil);
		return;
	}

	u.Password = "";
	lib.SendResponse(w, http.StatusOK, "Login successfully", u, nil, &c);
}
