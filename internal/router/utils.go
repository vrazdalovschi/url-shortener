package router

import (
	"encoding/json"
	"net/http"
)

func RespondJson(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

func RespondOkJson(w http.ResponseWriter, body interface{}) {
	RespondJson(w, http.StatusOK, body)
}
