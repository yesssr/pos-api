package router

import (
	"pos-api/internal/handler"
	"pos-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func New(h *handler.Handler) chi.Router {
	r := chi.NewRouter();

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			// Docs
			r.Mount("/docs", DocsRouter());
			r.Post("/auth/login", h.Auth.Login);
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(middleware.TokenAuth));
			r.Use(middleware.Auth);


			// Non-admin subgroup
			r.Group(func(r chi.Router) {
				r.Mount("/customers", CustomerRouter(h.Customer));
			});

			// Admin subgroup
			r.Group(func(r chi.Router) {
				r.Use(middleware.IsAdmin);
				r.Mount("/users", UserRouter(h.User));
				r.Mount("/products", ProductRouter(h.Product));
			});

		});
	});
	return r;
}
