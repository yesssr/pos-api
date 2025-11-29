package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"pos-api/internal/handler"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"
)

func ProductRouter(h *handler.ProductHandler) http.Handler {
	r := chi.NewRouter();
 	allowedCols := map[string]bool{
   	"name":   		true,
    "price":      true,
    "stock":      true,
    "created_at": true,
    "updated_at": true,
  }

	r.Post("/", h.CreateProduct);
	r.With(lib.Paginate, middleware.QueryCtx(allowedCols)).Get("/", h.ListProducts);
	r.Route("/{id}", func(r chi.Router) {
		r.Use(middleware.IdCtx);
		r.Get("/", h.GetProduct);
		r.Put("/", h.UpdateProduct);
		r.Delete("/", h.DeleteProduct);
	});
	r.Get("/count", h.GetTotalProduct);
	return r;
}
