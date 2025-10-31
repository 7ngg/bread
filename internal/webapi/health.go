package webapi

import (
	"net/http"
)

func alive(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("alive"))
}
