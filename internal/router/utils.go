package router

import (
	"encoding/json"
	"github.com/vrazdalovschi/url-shortener/internal/domain"
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

func RespondError(w http.ResponseWriter, err domain.Error) {
	if err.ErrorCode == 0 {
		err.ErrorCode = 500
	}
	RespondJson(w, err.ErrorCode, err)
}
