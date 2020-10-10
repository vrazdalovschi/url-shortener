package main

import (
	"fmt"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	serivcemiddleware "github.com/vrazdalovschi/url-shortener/internal/middleware/serivce"
	"github.com/vrazdalovschi/url-shortener/internal/repository"
	"github.com/vrazdalovschi/url-shortener/internal/repository/postgres"
	"github.com/vrazdalovschi/url-shortener/internal/router"
	"github.com/vrazdalovschi/url-shortener/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"

	"net/http"
)

func main() {
	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "", flag.ExitOnError)
	var (
		httpAddr   = fs.String("HTTP_ADDR", ":8080", "HTTP endpoint to use for endpoints")
		dbHost     = fs.String("DB_HOST", "localhost", "DB Host")
		dbPort     = fs.String("DB_PORT", "5432", "DB Port")
		dbUser     = fs.String("DB_USER", "url-shortener", "DB Username")
		dbPassword = fs.String("DB_PASSWORD", "root", "DB Password")
		dbName     = fs.String("DB_NAME", "shortener", "DB name")
	)
	_ = fs.Parse(os.Args[1:])

	labelNames := []string{"method", "error"}
	counterOpts := stdprometheus.CounterOpts{
		Namespace: "url_shortener",
		Subsystem: "shortener",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}
	summaryOpts := stdprometheus.SummaryOpts{
		Namespace: "url_shortener",
		Subsystem: "shortener",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}
	requestCounter := prometheus.NewCounterVec(counterOpts, labelNames)
	requestLatencySummary := prometheus.NewSummaryVec(summaryOpts, labelNames)
	prometheus.MustRegister(requestCounter, requestLatencySummary)

	configuration := repository.Configuration{
		Host:     *dbHost,
		Port:     *dbPort,
		User:     *dbUser,
		Password: *dbPassword,
		DbName:   *dbName,
	}
	st, err := postgres.NewRepository(configuration)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	var svc service.Service
	{
		svc = service.NewService(st)
		svc = serivcemiddleware.NewLogging()(svc)
		svc = serivcemiddleware.NewMetrics(requestCounter, requestLatencySummary)(svc)
	}
	r := router.New(svc)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		log.Println("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, r)
	}()

	log.Fatal("exit", <-errs)
}
