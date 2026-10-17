-- +goose Up
CREATE TABLE otp_challenges (
    id uuid PRIMARY KEY,
    purpose text NOT NULL,
    target text NOT NULL,
    user_id uuid REFERENCES users (id),
    code_hash bytea NOT NULL,
    attempts smallint NOT NULL DEFAULT 0,
    expires_at timestamptz NOT NULL,
    consumed_at timestamptz,
    request_ip inet NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT otp_challenges_purpose_chk CHECK (purpose IN ('LOGIN', 'PHONE_VERIFY')),
    CONSTRAINT otp_challenges_attempts_chk CHECK (attempts BETWEEN 0 AND 5)
);

CREATE INDEX otp_challenges_target_created_idx ON otp_challenges(target, created_at);
CREATE INDEX otp_challenges_ip_created_idx ON otp_challenges(request_ip, created_at);

-- +goose Down
DROP TABLE IF EXISTS otp_challenges;