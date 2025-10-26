package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	chi := chi.NewRouter()

	chi.Use(middleware.Logger)

	http.ListenAndServe(":42069", chi)
}
