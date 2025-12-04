package router

import (
	"net/http"
	"pos-api/internal/handler"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func CustomerRouter(h *handler.CustomerHandler) http.Handler {
	r := chi.NewRouter();
 	allowedCols := map[string]bool{
   	"name":   true,
    "created_at": true,
  }

	r.Post("/", h.CreateCustomer);
	r.With(middleware.QueryCtx(allowedCols), lib.Paginate).Get("/", h.ListCustomers);
	r.Route("/{id}", func(r chi.Router) {
		r.Use(middleware.IdCtx);
		r.Get("/", h.GetCustomer);
		r.Put("/", h.UpdateCustomer);
		r.Delete("/", h.DeleteCustomer);
	});
	r.Get("/count", h.GetTotalCustomer);
	return r;
}
