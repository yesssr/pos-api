package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func DocsRouter() http.Handler {
	r := chi.NewRouter();
	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte(
				`<!DOCTYPE html>
					<html lang="en">
						<head>
							<meta charset="UTF-8" />
							<title>POS API Documentation</title>
							<link rel="icon" href="data:,">
							<style>
							  html,body { margin:0; padding:0; height:100%; }
							  #redoc { height:100vh; }
							</style>
						</head>
						<body>
							<div id="redoc"></div>
							<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
							<script>
							  Redoc.init('/api/v1/docs/openapi.yaml', {}, document.getElementById('redoc'));
							</script>
						</body>
					</html>`,
			));
		});


		// Serve OpenAPI YAML files
		r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/vnd.oai.openapi+yaml; charset=utf-8")
			http.ServeFile(w, r, "docs/openapi.yaml");
		});
		r.Get("/auth.yaml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml; charset=utf-8");
			http.ServeFile(w, r, "docs/auth.yaml");
		});
		r.Get("/users.yaml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml; charset=utf-8");
			http.ServeFile(w, r, "docs/users.yaml");
		});
		r.Get("/customers.yaml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml; charset=utf-8");
			http.ServeFile(w, r, "docs/customers.yaml");
		});
		r.Get("/products.yaml", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/yaml; charset=utf-8");
			http.ServeFile(w, r, "docs/products.yaml");
		});
	});

	return r;
}
