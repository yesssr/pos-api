package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"pos-api/internal/handler"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"
)

func UserRouter(h *handler.UserHandler) http.Handler {
  r := chi.NewRouter();
  r.Post("/", h.CreateUser);
  r.With(lib.Paginate).Get("/", h.ListUsers);
  r.Route("/{id}", func(r chi.Router) {
  	r.Use(middleware.IdCtx)
   	r.Get("/", h.GetUser);
    r.Put("/", h.UpdateUser);
    r.Delete("/", h.DeleteUser);
  })
  return r
}
