package web

import (
	"database/sql"
	"net/http"

	"github.com/7ngg/bread/internal/db"
	"github.com/7ngg/bread/internal/services"
	"github.com/7ngg/bread/internal/web/api"
	"github.com/7ngg/bread/internal/web/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

type WebApp struct {
	db             *db.Queries
	cache          *redis.Client
	router         *chi.Mux
	productService *services.ProductService
	basketService  *services.BasketService
}

func NewWebApp(dbConn *sql.DB, redisClient *redis.Client) *WebApp {
	queries := db.New(dbConn)
	productService := services.NewProductService(queries)
	basketService := services.NewBasketService(redisClient)
	app := &WebApp{
		db:             queries,
		cache:          redisClient,
		router:         chi.NewRouter(),
		productService: productService,
		basketService:  basketService,
	}

	app.router.Use(middleware.Logger)

	// Health check endpoint
	healthHandler := api.NewHealthHandler(redisClient, dbConn)
	app.router.Get("/api/_health", healthHandler.HealthCheck)
	app.router.Get("/api/_health/ready", healthHandler.Readiness)
	app.router.Get("/api/_health/alive", healthHandler.Liveness)

	// Products endpoints
	productHandler := api.NewProductsHandler(productService)
	app.router.Get("/api/products", productHandler.GetAllProductsHandler)

	// Pages
	pagesHadnler := pages.NewPagesHandler(productService, basketService)
	fs := http.FileServer(http.Dir("static"))
	app.router.Handle("/static/*", http.StripPrefix("/static/", fs))
	app.router.Get("/", pagesHadnler.RenderIndex)
	app.router.Post("/basket", pagesHadnler.RenderPlus)

	return app
}

func (handler *WebApp) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, handler.router)
}

