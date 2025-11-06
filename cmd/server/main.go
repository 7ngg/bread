package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/7ngg/bread/internal/cache"
	"github.com/7ngg/bread/internal/config"
	"github.com/7ngg/bread/internal/web"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoad()

	dbConn := initializeDbConnection(cfg.DbConn)
	defer dbConn.Close()

	redisClient := cache.NewRedisConnection(&cfg.Redis)
	defer redisClient.Close()

	app := web.NewWebApp(dbConn, redisClient)

	log.Fatal(app.ListenAndServe(fmt.Sprintf(":%d", cfg.Port)))
}

func initializeDbConnection(connectionString string) *sql.DB {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return conn
}

