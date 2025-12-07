package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/7ngg/bread/internal/cache"
	"github.com/7ngg/bread/internal/config"
	"github.com/7ngg/bread/internal/web"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	cfg := config.MustLoad()

	dbConn := initializeDbConnection(cfg.DbConn)
	defer dbConn.Close()

	redisClient := cache.NewRedisConnection(&cfg.Redis)
	defer redisClient.Close()

	app := web.NewWebApp(dbConn, redisClient, logger)

	log.Fatal(app.ListenAndServe(fmt.Sprintf(":%d", cfg.Port)))
}

func initializeDbConnection(connectionString string) *pgxpool.Pool {
	conn, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return conn
}
