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
  allowedCols := map[string]bool{
   	"username":   true,
		"role":       true,
    "created_at": true,
    "updated_at": true,
  }

  r.Post("/", h.CreateUser);
  r.With(middleware.QueryCtx(allowedCols), lib.Paginate).Get("/", h.ListUsers);
  r.Route("/{id}", func(r chi.Router) {
  	r.Use(middleware.IdCtx);
   	r.Get("/", h.GetUser);
    r.Put("/", h.UpdateUser);
    r.Delete("/", h.DeleteUser);
  });
  r.Get("/count", h.GetTotalUser);
  return r;
}
