package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"pos-api/internal/config"
	"pos-api/internal/handler"
	"pos-api/internal/router"
	"pos-api/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system env")
  }
  ctx := context.Background();
  db, err := config.NewPool(ctx);
  if err != nil {
		panic(err);
	}
	q := store.New(db);
	h := handler.New(q);
	r := router.New(h);

	port := os.Getenv("PORT")

	fmt.Println("Starting server on port", port)
	http.ListenAndServe(":"+port, r);
}
