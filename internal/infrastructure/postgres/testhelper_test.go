package postgres_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../../.env")

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		panic("TEST_DATABASE_URL environment variable is not set")
	}

	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		panic("failed to open sql.DB for migrations: " + err.Error())
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic("failed to set goose dialect: " + err.Error())
	}

	if err := goose.Up(sqlDB, "../../../db/migrations"); err != nil {
		panic("failed to run goose migrations: " + err.Error())
	}

	if err := sqlDB.Close(); err != nil {
		panic("failed to close sql.DB after migrations: " + err.Error())
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic("failed to create pgxpool: " + err.Error())
	}
	testPool = pool

	code := m.Run()

	pool.Close()
	os.Exit(code)
}

func cleanDB(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		"TRUNCATE TABLE credit_ledger, apps RESTART IDENTITY CASCADE",
	)
	if err != nil {
		t.Fatalf("cleanDB: failed to truncate tables: %v", err)
	}
}
