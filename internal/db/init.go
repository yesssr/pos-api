package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"pos-api/internal/lib"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system env")
  }

  cString := os.Getenv("DATABASE_URL");
  pool, err := sql.Open("postgres", cString);
  if err != nil {
    log.Fatalf("pool init: %v", err)
  }
  defer pool.Close()

  driver, err := postgres.WithInstance(pool, &postgres.Config{})
  if err != nil {
    log.Fatalf("driver: %v", err)
  }

  m, err := migrate.NewWithDatabaseInstance(
    "file://internal/db/migrations",
    "postgres",
    driver,
  )
  if err != nil {
    log.Fatalf("migrate init: %v", err)
  }

  if err := m.Down(); err != nil && err.Error() != "no change" {
    log.Fatalf("migrate down: %v", err)
  }

  if err := m.Up(); err != nil && err.Error() != "no change" {
    log.Fatalf("migrate up: %v", err)
  }

  runSeeds(pool)
}

func runSeeds(pool *sql.DB) {
  p, _ := lib.HashPassword("00000000")
  q := "INSERT INTO users (username, password, role) VALUES ('admin', $1, 'admin')"
  _, err := pool.Exec(q, p)
  if err != nil {
    log.Fatalf("seed admin: %v", err)
  }
}
