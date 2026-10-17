package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/codecommander03/hestia/backend/internal/db"
)

type CreateUserInput struct {
	Body struct {
		Name  string `json:"name" minLength:"1" maxLength:"100"`
		Email string `json:"email" format:"email"`
	}
}

type GetUserInput struct {
	ID string `path:"id" format:"uuid"`
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

// newUserOutput converts a database row into the shape the API returns.
func newUserOutput(user db.User) *UserOutput {
	resp := &UserOutput{}
	resp.Body.ID = user.ID.String()
	resp.Body.Name = user.Name.String
	resp.Body.Email = user.Email.String
	resp.Body.Status = user.Status
	resp.Body.CreatedAt = user.CreatedAt.Time
	return resp
}

// Registers Users the /users endpoints to api.
func RegisterUsers(api huma.API, queries *db.Queries) {
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
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return nil, huma.Error409Conflict("email already registered")
			}
			log.Printf("create user: %v", err)
			return nil, huma.Error500InternalServerError("error creating the user")
		}

		return newUserOutput(user), nil
	})

	huma.Get(api, "/users/{id}", func(ctx context.Context, input *GetUserInput) (*UserOutput, error) {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid user id")
		}
		user, err := queries.GetUserByID(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, huma.Error404NotFound("user not found")
			}
			log.Printf("get user: %v", err)
			return nil, huma.Error500InternalServerError("could not get user")
		}

		return newUserOutput(user), nil
	})
}
