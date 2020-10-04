package main

import (
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"os"
	"os/signal"
	"syscall"

	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api", IsNotImplementedHandler)
	r.HandleFunc("/api/{shortedUrlId}", IsNotImplementedHandler)

	fs := http.FileServer(http.Dir("./api/"))
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swaggerui/swagger.yaml"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))
	http.Handle("/", r)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		log.Println("transport", "HTTP", "addr", "localhost:8085")
		errs <- http.ListenAndServe("localhost:8085", r)
	}()

	log.Println("exit", <-errs)
}

func IsNotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	if _, err := w.Write([]byte("It's not implemented")); err != nil {
		log.Printf("Got %v for request: %s", err, r.URL.Path)
	}
}
