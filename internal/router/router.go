package router

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vrazdalovschi/url-shortener/internal/service"
	"log"
	"net/http"
)

func New(_ service.Service, host string) *mux.Router {
	r := mux.NewRouter()

	r.Methods("GET").Path("/{shortedUrlId}").HandlerFunc(IsNotImplementedHandler)
	r.Methods("POST").Path("/api").HandlerFunc(IsNotImplementedHandler)
	r.Methods("GET").Path("/api/{shortedUrlId}").HandlerFunc(IsNotImplementedHandler)
	r.Methods("DELETE").Path("/api/{shortedUrlId}").HandlerFunc(IsNotImplementedHandler)

	fs := http.FileServer(http.Dir("./api/"))
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(host+"/swaggerui/swagger.yaml"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))
	http.Handle("/", r)

	return r
}

func IsNotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	if _, err := w.Write([]byte("It's not implemented")); err != nil {
		log.Printf("Got %v for request: %s", err, r.URL.Path)
	}
}
