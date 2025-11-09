package web

import (
	"database/sql"
	"html/template"
	"io"
	"net/http"

	"github.com/7ngg/bread/internal/db"
	"github.com/7ngg/bread/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

type WebApp struct {
	templates      *template.Template
	db             *db.Queries
	cache          *redis.Client
	router         *chi.Mux
	productService *services.ProductService
	basketService  *services.BasketService
}

func NewWebApp(dbConn *sql.DB, redisClient *redis.Client) *WebApp {
	queries := db.New(dbConn)
	app := &WebApp{
		templates:      NewTemplates(),
		db:             queries,
		cache:          redisClient,
		router:         chi.NewRouter(),
		productService: services.NewProductService(queries),
		basketService:  services.NewBasketService(redisClient),
	}

	healthHandler := NewHealthHandler(redisClient, dbConn)

	app.router.Use(middleware.Logger)

	// Health check endpoint
	app.router.Get("/api/_health", healthHandler.HealthCheck)
	app.router.Get("/api/_health/ready", healthHandler.Readiness)
	app.router.Get("/api/_health/alive", healthHandler.Liveness)

	// Products endpoints
	app.router.Get("/api/products", app.GetAllProductsHandler)

	// Pages
	fs := http.FileServer(http.Dir("static"))
	app.router.Handle("/static/*", http.StripPrefix("/static/", fs))
	app.router.Get("/", app.RenderIndex)
	app.router.Post("/basket", app.RenderAddItemToBasket)

	return app
}

func (handler *WebApp) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, handler.router)
}

func (handler *WebApp) Render(w io.Writer, name string, data any) error {
	return handler.templates.ExecuteTemplate(w, name, data)
}
