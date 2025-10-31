package webapi

import (
	"database/sql"

	"github.com/7ngg/bread/internal/config"
	"github.com/7ngg/bread/internal/db"
	"github.com/7ngg/bread/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"
)

type App struct {
	db *db.Queries
	productService services.ProductService
}

func NewApp(cfg *config.AppConfig) *App {
	dbConn, err := sql.Open("postgres", cfg.DbPath)
	if err != nil {
		panic("unable to connect to database")
	}

	db := db.New(dbConn)

	return &App{
		db: db,
		productService: services.ProductService{
			ProductsGetter: db,
		},
	}
}

func NewRouter(cfg *config.AppConfig) *chi.Mux {
	r := chi.NewRouter()
	app := NewApp(cfg)

	r.Use(middleware.Logger)

	// Health check endpoint
	r.Get("/api/_health/alive", alive)

	// Products endpoints
	r.Get("/api/products", app.GetAllProductsHandler)

	return r;
}
