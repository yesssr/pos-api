package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"pos-api/internal/handler"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"
)

func ProductRouter(h *handler.ProductHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.CreateProduct)
	r.With(lib.Paginate, middleware.QueryCtx).Get("/", h.ListProducts)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(middleware.IdCtx)
		r.Get("/", h.GetProduct)
		r.Put("/", h.UpdateProduct)
		r.Delete("/", h.DeleteProduct)
	})
	r.Get("/count", h.GetTotalProduct)
	return r
}
