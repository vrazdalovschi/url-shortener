package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vrazdalovschi/url-shortener/internal/domain"
	"github.com/vrazdalovschi/url-shortener/internal/service"
	"github.com/vrazdalovschi/url-shortener/internal/stackerr"
	"net/http"
)

func New(svc service.Service) *mux.Router {
	r := mux.NewRouter()

	r.Methods("GET").Path("/{shortenedId}").HandlerFunc(Redirect(svc))
	r.Methods("POST").Path("/api").HandlerFunc(CreateShortenedId(svc))
	r.Methods("GET").Path("/api/{shortenedId}").HandlerFunc(Describe(svc))
	r.Methods("DELETE").Path("/api/{shortenedId}").HandlerFunc(Delete(svc))
	r.Methods("GET").Path("/api/stats/{shortenedId}").HandlerFunc(Stats(svc))

	fs := http.FileServer(http.Dir("./api/"))
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swaggerui/swagger.yaml"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))
	r.Handle("/internal/metrics", promhttp.Handler())
	http.Handle("/", r)

	return r
}

func CreateShortenedId(svc service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var createRequest domain.CreateShortId
		if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
			e := domain.Error{
				Message:   fmt.Sprintf("Invalid request body. Err: %v", stackerr.Wrap(err)),
				ErrorCode: 400,
			}
			RespondError(w, e)
			return
		}
		res, err := svc.CreateShort(r.Context(), createRequest.ApiKey, createRequest.OriginalURL, createRequest.ExpiryDate)
		if err != nil {
			if wrapped, ok := err.(stackerr.HandledError); ok {
				if e, ok := wrapped.Inner().(domain.Error); ok {
					RespondError(w, e)
					return
				}
			}
			e := domain.Error{Message: fmt.Sprintf("Unexpected error: %v", err), ErrorCode: 500}
			RespondError(w, e)
			return
		}
		RespondOkJson(w, res)
	}
}

func Redirect(svc service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["shortenedId"]
		res, err := svc.GetOriginalUrl(r.Context(), id)
		if err != nil {
			e := domain.Error{Message: fmt.Sprintf("shortenedId %s not found", id), ErrorCode: 404}
			RespondError(w, e)
			return
		}
		//Ignore error here, it need just for logging middleware. Increment shouldn't affect redirect logic
		_ = svc.IncrementStats(r.Context(), id)
		http.Redirect(w, r, res, http.StatusFound)
	}
}

func Describe(svc service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["shortenedId"]
		res, err := svc.Describe(r.Context(), id)
		if err != nil {
			e := domain.Error{Message: fmt.Sprintf("shortenedId %s not found", id), ErrorCode: 404}
			RespondError(w, e)
			return
		}
		RespondOkJson(w, res)
	}
}

func Stats(svc service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["shortenedId"]
		res, err := svc.Stats(r.Context(), id)
		if err != nil {
			e := domain.Error{Message: fmt.Sprintf("shortenedId %s not found", id), ErrorCode: 404}
			RespondError(w, e)
			return
		}
		RespondOkJson(w, res)
	}
}

func Delete(svc service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := mux.Vars(r)["shortenedId"]
		err := svc.Delete(r.Context(), id)
		if err != nil {
			e := domain.Error{Message: fmt.Sprintf("shortenedId %s not found", id), ErrorCode: 404}
			RespondError(w, e)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
