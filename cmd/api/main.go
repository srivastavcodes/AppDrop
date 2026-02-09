package main

import (
	"appdrop/internal/data"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

import _ "github.com/lib/pq"

func init() {
	if err := godotenv.Load(".envrc"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// todo: use flags for configuration later
	cfg := config{
		port: 8080, db: struct {
			dsn          string
			maxOpenConns int
			maxIdleTime  time.Duration
			maxIdleConns int
		}{
			dsn:          os.Getenv("APP_DROP_DSN"),
			maxOpenConns: 25,
			maxIdleTime:  90 * time.Second,
			maxIdleConns: 25,
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDb(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err.Error())
		}
	}()
	logger.Info("database connection established")

	b := &backend{
		logger: logger,
		conf:   cfg,
		models: data.NewModels(db),
	}
	if err = b.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDb(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, fmt.Errorf("couldn't open database: %w", err)
	}
	db.SetConnMaxLifetime(cfg.db.maxIdleTime)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("couldn't connect to database: %w", err)
	}
	return db, nil
}
