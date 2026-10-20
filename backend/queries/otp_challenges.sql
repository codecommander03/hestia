-- name: CreateOtpChallenge :one
INSERT INTO otp_challenges (id, purpose, target, user_id, code_hash, expires_at, request_ip)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: CountOtpChallengesByTargetSince :one
SELECT count(*) FROM otp_challenges
WHERE target = sqlc.arg(target) AND created_at > sqlc.arg(since);

-- name: CountOtpChallengesByIPSince :one
SELECT count(*) FROM otp_challenges
WHERE request_ip = sqlc.arg(request_ip) AND created_at > sqlc.arg(since);

-- name: GetLatestLiveOtpChallenge :one
SELECT * FROM otp_challenges
WHERE target = $1 AND purpose = $2
    AND consumed_at IS NULL
    and expires_at > NOW()
    and attempts < 5
ORDER BY created_at DESC
LIMIT 1;

-- name: IncrementOtpChallengeAttempts :one
UPDATE otp_challenges
SET attempts = attempts + 1, updated_at = NOW()
WHERE id = $1 AND attempts < 5
RETURNING attempts;

-- name: ConsumeOtpChallenge :execrows
UPDATE otp_challenges
SET consumed_at = now(), updated_at = now()
WHERE id = $1 AND consumed_at IS NULL;