package web

import (
	"net/http"

	"github.com/7ngg/bread/internal/db"
	"github.com/7ngg/bread/internal/services"
	"github.com/7ngg/bread/internal/web/api"
	"github.com/7ngg/bread/internal/web/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type WebApp struct {
	router         *chi.Mux
}

func NewWebApp(dbConn *pgxpool.Pool, redisClient *redis.Client) *WebApp {
	queries := db.New(dbConn)
	productService := services.NewProductService(queries)
	basketService := services.NewBasketService(redisClient)
	userService := services.NewUserService(queries)
	orderService := services.NewOrderService(queries, queries, queries, userService)
	app := &WebApp{
		router:         chi.NewRouter(),
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
	pagesHadnler := pages.NewPagesHandler(productService, basketService, orderService)
	fs := http.FileServer(http.Dir("static"))
	app.router.Handle("/static/*", http.StripPrefix("/static/", fs))
	app.router.Get("/", pagesHadnler.RenderIndex)
	app.router.Get("/checkout", pagesHadnler.RenderCheckout)
	app.router.Post("/basket", pagesHadnler.RenderPlus)
	app.router.Post("/checkout", pagesHadnler.RenderOrder)

	return app
}

func (handler *WebApp) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, handler.router)
}
