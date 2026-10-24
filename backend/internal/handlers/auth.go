package handlers

import (
	"context"
	"log"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/codecommander03/hestia/backend/internal/db"
	"github.com/codecommander03/hestia/backend/internal/otp"
)

const (
	otpTTL         = 10 * time.Minute
	otpMaxPerHour  = 5
	otpLimitWindow = time.Hour
)

type StartLoginInput struct {
	Body struct {
		Email string `json:"email" format:"email"`
	}
	ip netip.Addr
}

// Resolve runs after Huma has parsed the request and saves the caller's IP address.
func (i *StartLoginInput) Resolve(ctx huma.Context) []error {
	addrPort, err := netip.ParseAddrPort(ctx.RemoteAddr())
	if err != nil {
		return []error{&huma.ErrorDetail{Message: "could not read client address"}}
	}
	i.ip = addrPort.Addr()
	return nil
}

type StartLoginOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

// RegisterAuth registers the /auth endpoints to api.
func RegisterAuth(api huma.API, queries *db.Queries, otpSecret []byte) {
	huma.Register(api, huma.Operation{
		OperationID:   "start-login",
		Method:        http.MethodPost,
		Path:          "/auth/login/start",
		Summary:       "Emails a 6-digit login code",
		DefaultStatus: http.StatusAccepted,
	}, func(ctx context.Context, input *StartLoginInput) (*StartLoginOutput, error) {
		email := strings.ToLower(strings.TrimSpace(input.Body.Email))
		since := pgtype.Timestamptz{Time: time.Now().Add(-otpLimitWindow), Valid: true}
		sentToEmail, err := queries.CountOtpChallengesByTargetSince(ctx, db.CountOtpChallengesByTargetSinceParams{
			Target: email,
			Since:  since,
		})
		if err != nil {
			log.Printf("count otps by target: %v", err)
			return nil, huma.Error500InternalServerError("could not send code")
		}
		sentFromIP, err := queries.CountOtpChallengesByIPSince(ctx, db.CountOtpChallengesByIPSinceParams{
			RequestIp: input.ip,
			Since:     since,
		})
		if err != nil {
			log.Printf("count otps by ip: %v", err)
			return nil, huma.Error500InternalServerError("could not send code")
		}
		if sentToEmail >= otpMaxPerHour || sentFromIP >= otpMaxPerHour {
			return nil, huma.Error429TooManyRequests("too many codes requested. try again later")
		}

		code, err := otp.Generate()
		if err != nil {
			return nil, err
		}
		id, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}

		_, err = queries.CreateOtpChallenge(ctx, db.CreateOtpChallengeParams{
			ID:        id,
			Purpose:   "LOGIN",
			Target:    email,
			CodeHash:  otp.Hash(otpSecret, code),
			ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(otpTTL), Valid: true},
			RequestIp: input.ip,
		})
		if err != nil {
			log.Printf("create otp challenge: %v, err")
			return nil, huma.Error500InternalServerError("could not send code")
		}

		// Dev provider: print the code instead of emailing it
		log.Printf("[dev] login code for %s: %s", email, code)

		resp := &StartLoginOutput{}
		resp.Body.Message = "if the address is valid, a code is on its way"
		return resp, nil
	})
}
