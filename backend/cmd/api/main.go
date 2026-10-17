package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/codecommander03/hestia/backend/internal/db"
)

type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok"`
	}
}

type CreateUserInput struct {
	Body struct {
		Name  string `json:"name" minLength:"1" maxlength:"100"`
		Email string `json:"email" format: "email"`
	}
}

type UserOutput struct {
	Body struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
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

	queries := db.New(pool)

	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("Hestia API", "0.1.0"))

	huma.Get(api, "/health", func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "ok"
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "/users",
		Summary:       "Creates user (temporary, replaced by email-code sign up later)",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *CreateUserInput) (*UserOutput, error) {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}

		user, err := queries.CreateUser(ctx, db.CreateUserParams{
			ID:    id,
			Name:  pgtype.Text{String: input.Body.Name, Valid: true},
			Email: pgtype.Text{String: input.Body.Email, Valid: true},
		})
		if err != nil {
			log.Printf("create user: %v", err)
			return nil, huma.Error500InternalServerError("could not create user")
		}

		resp := &UserOutput{}
		resp.Body.ID = user.ID.String()
		resp.Body.Name = user.Name.String
		resp.Body.Email = user.Email.String
		resp.Body.Status = user.Status
		resp.Body.CreatedAt = user.CreatedAt.Time
		return resp, nil
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
