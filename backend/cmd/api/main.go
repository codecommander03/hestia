package main

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok"`
	}
}

func main() {
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("Hestia API", "0.1.0"))

	huma.Get(api, "/health", func(ctx context.Context, input *struct{})  (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = "ok"
		return resp, nil
	})

	http.ListenAndServe(":8080", mux)
}
