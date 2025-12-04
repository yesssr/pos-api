package router

import (
	"pos-api/internal/handler"
	"pos-api/internal/lib"
	"pos-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

func New(h *handler.Handler) chi.Router {
	r := chi.NewRouter();
	r.Use(cors.Handler(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}));

	allowedProductsCols := map[string]bool{
   	"name":   		true,
    "price":      true,
    "stock":      true,
    "created_at": true,
    "updated_at": true,
  }

	allowedPeriod := map[string]bool{
    "day": true,
    "week": true,
    "month": true,
    "year": true,
  }

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			// Docs
			r.Mount("/docs", DocsRouter());
			r.Post("/auth/login", h.Auth.Login);
			r.Post("/webhooks/xendit", h.Transaction.WebHookXendit);
		});

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(middleware.TokenAuth));
			r.Use(middleware.Auth);


			// Non-admin subgroup
			r.Group(func(r chi.Router) {
				r.Mount("/customers", CustomerRouter(h.Customer));
				r.With(lib.Paginate, middleware.QueryCtx(allowedProductsCols, map[string]bool{})).
					Get("/products-active", h.Product.ListProductsActive);
				r.Post("/cashier/transactions", h.Transaction.CreateTransaction);
			});

			// Admin subgroup
			r.Group(func(r chi.Router) {
				r.Use(middleware.IsAdmin);
				r.With(lib.Paginate, middleware.QueryCtx(map[string]bool{}, map[string]bool{})).
					Get("/transactions", h.Transaction.ListTransactions);
				r.With(lib.Paginate, middleware.QueryCtx(map[string]bool{}, allowedPeriod)).
					Get("/list-sales", h.Transaction.SalesByPeriods);
				r.Mount("/users", UserRouter(h.User));
				r.Mount("/products", ProductRouter(h.Product, allowedProductsCols));
			});

		});
	});
	return r;
}
