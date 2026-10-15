-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    name text,
    email citext UNIQUE,
    phone text UNIQUE,
    phone_verified_at timestamptz,
    trust_tier smallint NOT NULL DEFAULT 0,
    status text NOT NULL DEFAULT 'ACTIVE',
    suspended_reason text,
    carrier_declaration_at timestamptz,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT users_trust_tier_chk CHECK (trust_tier IN(0, 1)),
    CONSTRAINT users_status_chk CHECK (status IN ('ACTIVE', 'SUSPENDED', 'DELETED')),
    CONSTRAINT users_phone_e164_chk CHECK (phone ~ '^\+[1-9][0-9]{7,14}$'),
    CONSTRAINT users_suspended_reason_chk CHECK ((status = 'SUSPENDED') = (suspended_reason IS NOT NULL)),
    CONSTRAINT users_deleted_at_chk CHECK ((status = 'DELETED') = (deleted_at IS NOT NULL)),
    CONSTRAINT users_email_required_chk CHECK (status = 'DELETED' OR email IS NOT NULL)
);

-- +goose Down
DROP TABLE IF EXISTS users;
