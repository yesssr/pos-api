package router

import (
	"pos-api/internal/handler"
	"pos-api/internal/lib"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func New(h *handler.Handler) chi.Router {
	r := chi.NewRouter();

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/auth/login", h.Auth.Login);
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(lib.TokenAuth));
			r.Use(jwtauth.Authenticator(lib.TokenAuth));

			r.Route("/users", func(r chi.Router) {
				r.Post("/", h.User.CreateUser);
			});
		})
	});
	return r;
}
