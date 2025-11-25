package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"pos-api/internal/configuration"
	"pos-api/internal/handler"
	"pos-api/internal/lib"
	"pos-api/internal/router"
	"pos-api/internal/service"
	"pos-api/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system env")
  }
  lib.InitJWT();
  ctx := context.Background();

  db, err := configuration.NewPool(ctx);
  if err != nil {
		panic(err);
	}

  awsClient, err := configuration.NewAwsClient();
  if err != nil {
		panic(err);
	}

	bucket := os.Getenv("BUCKET_NAME");
	q := store.New(db);
	s := service.New(q, awsClient, bucket);
	h := handler.New(s);
	r := router.New(h);

	port := os.Getenv("PORT")

	fmt.Println("Starting server on port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		fmt.Println("Server error:", err);
	}
}
