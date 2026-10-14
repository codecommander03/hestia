package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok"`
	}
}

func main() {
	_ = godotenv.Load()

	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Create db pool: %v", err)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("ping db: %v", err)
	}
	log.Println("connected to database")

	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("Hestia API", "0.1.0"))

	huma.Get(api, "/health", func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "ok"
		return resp, nil
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
