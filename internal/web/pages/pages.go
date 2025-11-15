package pages

import (
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/7ngg/bread/internal/lib"
	"github.com/7ngg/bread/internal/services"
)

type Templates struct {
	templates *template.Template
}

func NewTemplates() *template.Template {
	funcMap := template.FuncMap{
		"multiply": func(a float64, b int) float64 {
			return a * float64(b)
		},
	}

	return template.Must(template.New("").Funcs(funcMap).ParseGlob("views/*.html"))
}

type PagesHandler struct {
	templates      *template.Template
	productService *services.ProductService
	basketService  *services.BasketService
}

func NewPagesHandler(ps *services.ProductService, bs *services.BasketService) *PagesHandler {
	return &PagesHandler{
		templates:      NewTemplates(),
		productService: ps,
		basketService:  bs,
	}
}

func (handler *PagesHandler) Render(w io.Writer, name string, data any) error {
	return handler.templates.ExecuteTemplate(w, name, data)
}

func setSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    lib.RandomString(16),
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

func getSessionID(r *http.Request) string {
	c, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}
	return c.Value
}
