package lib

import (
	"encoding/json"
	"net/http"
)

func ResponseWithJson(w http.ResponseWriter, code int, payload any) (error) {
	json, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.WriteHeader(code)
	w.Write(json)

	return nil
}
