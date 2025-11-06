package web

import (
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
	productService services.ProductService
}

func NewWebApp(db *db.Queries, redisClient *redis.Client) *WebApp {
	app := &WebApp{
		templates:      NewTemplates(),
		db:             db,
		cache:          redisClient,
		router:         chi.NewRouter(),
		productService: *services.NewProductService(db),
	}

	app.router.Use(middleware.Logger)

	// Health check endpoint
	app.router.Get("/api/_health/alive", alive)

	// Products endpoints
	app.router.Get("/api/products", app.GetAllProductsHandler)

	// Pages
	fs := http.FileServer(http.Dir("static"))
	app.router.Handle("/static/*", http.StripPrefix("/static/", fs))
	app.router.Get("/", app.RenderIndex)

	return app
}

func (app *WebApp) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, app.router)
}

func (app *WebApp) Render(w io.Writer, name string, data any) error {
	return app.templates.ExecuteTemplate(w, name, data)
}
