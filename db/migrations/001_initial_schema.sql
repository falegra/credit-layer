-- +goose Up
CREATE TABLE apps (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    api_key     TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE credit_ledger (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id            UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    user_id           TEXT NOT NULL,
    amount            BIGINT NOT NULL,
    description       TEXT,
    idempotency_key   TEXT UNIQUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_credit_ledger_app_user ON credit_ledger(app_id, user_id);

-- +goose Down
DROP INDEX idx_credit_ledger_app_user;
DROP TABLE credit_ledger;
DROP TABLE apps;