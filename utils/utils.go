package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ParseJSON(r *http.Request, body any) error {
	if r.Body == nil {
		log.Fatal("Missing request body")
	}
	return json.NewDecoder(r.Body).Decode(body)
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(value)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

var Validate = validator.New()