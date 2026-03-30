package main

import (
	"context"
	"credit-layer/internal/application"
	"credit-layer/internal/infrastructure/postgres"
	httpHandlder "credit-layer/internal/interfaces/http"
	"database/sql"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL required")
	}

	if err := runMigrations(dbURL); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}

	log.Println("connected to database")

	appRepo := postgres.NewAppRepository(pool)
	creditLedgerRepo := postgres.NewCreditLedgerRepository(pool)

	appUseCase := application.NewAppUseCase(appRepo)
	creditUseCase := application.NewCreditLedgerUseCase(creditLedgerRepo)

	router := httpHandlder.NewRouter(appUseCase, creditUseCase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func runMigrations(dbURL string) error {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(db, "db/migrations")
}
