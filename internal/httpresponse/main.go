package httpresponse

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(struct {
		Data any `json:"data"`
	}{Data: v})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
