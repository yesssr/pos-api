package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"pos-api/internal/handler"
)

func UserRouter(h *handler.UserHandler) http.Handler {
  r := chi.NewRouter()
  r.Post("/", h.CreateUser)
  return r
}
