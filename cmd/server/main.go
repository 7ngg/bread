package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/7ngg/bread/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {
	cfg := config.MustLoad()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templates := NewTemplates()
		templates.Render(w, "index", cfg)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r))
}
