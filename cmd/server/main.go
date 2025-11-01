package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/7ngg/bread/internal/config"
	"github.com/7ngg/bread/internal/db"
	"github.com/7ngg/bread/internal/webapi"
	_ "github.com/lib/pq"
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

	r := webapi.NewRouter(cfg)

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templates := NewTemplates()
		conn, _ := sql.Open("postgres", cfg.DbPath)
		q := db.New(conn)
		products, _ := q.GetProducts(r.Context(), db.GetProductsParams{
			Limit:  10,
			Offset: 0,
		})
		opts := struct {
			Items []struct {
				Href string
				Text string
			}
			Products []db.Product
		}{
			Items: []struct {
				Href string
				Text string
			}{
				{
					Href: "/",
					Text: "Main",
				},
				{
					Href: "/",
					Text: "Catalog",
				},
			},
			Products: products,
		}
		err := templates.Render(w, "index.html", opts)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r))
}
