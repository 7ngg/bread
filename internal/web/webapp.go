package web

import (
	"net/http"

	"github.com/7ngg/bread/internal/db"
	"github.com/7ngg/bread/internal/services"
	"github.com/7ngg/bread/internal/web/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

type WebApp struct {
	templates      *pages.Templates
	db             *db.Queries
	cache          *redis.Client
	productService services.ProductService
	router         *chi.Mux
}

func NewWebApp(db *db.Queries, redisClient *redis.Client) *WebApp {
	app := &WebApp{
		templates: pages.NewTemplates(),
		db:        db,
		cache:     redisClient,
		router:    chi.NewRouter(),
		productService: services.ProductService{
			ProductsGetter: db,
		},
	}

	app.router.Use(middleware.Logger)

	// Health check endpoint
	app.router.Get("/api/_health/alive", alive)

	// Products endpoints
	app.router.Get("/api/products", app.GetAllProductsHandler)

	// Pages
	app.router.Get("/", app.templates.RenderIndex)

	return app
}

func (app *WebApp) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, app.router)
}
