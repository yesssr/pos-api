package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func main() {
	r := chi.NewRouter();

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		testResponse := Response{
			Success:    true,
			Message:    "API is working",
			StatusCode: 200,
		}
		w.WriteHeader(testResponse.StatusCode);
		w.Header().Set("Content-Type", "application/json");
		json.NewEncoder(w).Encode(testResponse);
	});

	fmt.Println("Starting server on port 8080");
	http.ListenAndServe(":8080", r);
}
