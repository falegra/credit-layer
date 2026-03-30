-- name: CreateApp :one
INSERT INTO apps (name, api_key)
VALUES ($1, $2)
RETURNING *;

-- name: GetAppByAPIKey :one
SELECT * FROM apps
WHERE api_key = $1;

-- name: AddCredits :one
INSERT INTO credit_ledger (app_id, user_id, amount, description, idempotency_key)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (idempotency_key) DO UPDATE
SET idempotency_key = EXCLUDED.idempotency_key
RETURNING *;

-- name: DeductCredits :one
INSERT INTO credit_ledger (app_id, user_id, amount, description, idempotency_key)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (idempotency_key) DO UPDATE
SET idempotency_key = EXCLUDED.idempotency_key
RETURNING *;

-- name: GetBalance :one
SELECT COALESCE(SUM(amount), 0)::bigint AS balance
FROM credit_ledger
WHERE app_id = $1 AND user_id = $2;

-- name: ExistsAppByName :one
SELECT EXISTS (
    SELECT 1 FROM apps WHERE name = $1
);