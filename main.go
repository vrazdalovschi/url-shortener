package main

import (
	"fmt"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/vrazdalovschi/url-shortener/internal/middleware/serivce"
	"github.com/vrazdalovschi/url-shortener/internal/router"
	"github.com/vrazdalovschi/url-shortener/internal/service"
	"github.com/vrazdalovschi/url-shortener/internal/storage/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"

	"net/http"
)

func main() {
	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "", flag.ExitOnError)
	var (
		httpAddr   = fs.String("HTTP_ADDR", "localhost:8080", "HTTP endpoint to use for endpoints")
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

	configuration := postgres.Configuration{
		Host:     *dbHost,
		Port:     *dbPort,
		User:     *dbUser,
		Password: *dbPassword,
		DbName:   *dbName,
	}
	st, err := postgres.New(configuration)
	if err != nil {
		log.Fatal(err)
	}

	var svc service.Service
	{
		svc = service.NewService(st)
		svc = serivcemiddleware.NewLogging()(svc)
		svc = serivcemiddleware.NewMetrics(requestCounter, requestLatencySummary)(svc)
	}
	r := router.New(svc, *httpAddr)

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
