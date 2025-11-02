package pages

import (
	"html/template"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func RegisterPages(r chi.Router) {
	templates := NewTemplates()
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", templates.RenderIndex)
}
