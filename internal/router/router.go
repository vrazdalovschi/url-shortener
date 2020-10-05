package router

import (
	"encoding/json"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vrazdalovschi/url-shortener/internal/service"
	"log"
	"net/http"
)

func New(svc service.Service, host string) *mux.Router {
	r := mux.NewRouter()

	r.Methods("GET").Path("/{shortedUrlId}").HandlerFunc(GetOriginalUrl(svc))
	r.Methods("POST").Path("/api").HandlerFunc(CreateShortenedId(svc, host))
	r.Methods("GET").Path("/api/{shortedUrlId}").HandlerFunc(Describe(svc))
	r.Methods("DELETE").Path("/api/{shortedUrlId}").HandlerFunc(Delete(svc))

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

func CreateShortenedId(svc service.Service, host string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var crReq createReq
		if err := json.NewDecoder(r.Body).Decode(&crReq); err != nil {
			log.Println(err)
			RespondJson(w, 400, err)
			return
		}
		res, err := svc.CreateShort(r.Context(), crReq.ApiKey, crReq.OriginalUrl, crReq.OriginalUrl)
		if err != nil {
			log.Println(err)
			RespondJson(w, 400, err)
			return
		}
		RespondOkJson(w, host+"/"+res)
	}
}

func GetOriginalUrl(svc service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["shortedUrlId"]
		if !ok || id == "" {
			RespondJson(w, 400, "bad request")
		}
		res, err := svc.GetOriginalUrl(r.Context(), id)
		if err != nil {
			log.Println(err)
			RespondJson(w, 400, err)
			return
		}
		http.Redirect(w, r, res, http.StatusFound)
	}
}

func Describe(svc service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["shortedUrlId"]
		if !ok || id == "" {
			RespondJson(w, 400, "bad request")
		}
		res, err := svc.Describe(r.Context(), id)
		if err != nil {
			log.Println(err)
			RespondJson(w, 400, err)
			return
		}
		RespondOkJson(w, res)
	}
}

func Delete(svc service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := mux.Vars(r)["shortedUrlId"]
		if !ok || id == "" {
			RespondJson(w, 400, "bad request")
		}
		err := svc.Delete(r.Context(), id)
		if err != nil {
			log.Println(err)
			RespondJson(w, 400, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

type createReq struct {
	ApiKey      string `json:"apiKey"`
	OriginalUrl string `json:"originalUrl"`
	ExpiryDate  string `json:"expiryDate"`
}
