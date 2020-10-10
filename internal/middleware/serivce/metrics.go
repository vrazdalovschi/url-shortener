package serivcemiddleware

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vrazdalovschi/url-shortener/internal/repository"
	"github.com/vrazdalovschi/url-shortener/internal/service"
	"time"
)

type metrics struct {
	rc   *prometheus.CounterVec
	rs   *prometheus.SummaryVec
	next service.Service
}

func NewMetrics(rc *prometheus.CounterVec, rs *prometheus.SummaryVec) Middleware {
	return func(s service.Service) service.Service {
		return &metrics{
			next: s,
			rc:   rc,
			rs:   rs,
		}
	}
}

func (m *metrics) CreateShort(ctx context.Context, apiKey, originalUrl, expiryDate string) (resp *repository.ShortenedIdResponse, err error) {
	defer func(begin time.Time) {
		elapsedTime := float64(time.Since(begin).Milliseconds())
		labels := prometheus.Labels{"method": "CreateShort", "error": fmt.Sprint(err != nil)}
		m.rc.With(labels).Add(1)
		m.rs.With(labels).Observe(elapsedTime)
	}(time.Now())
	return m.next.CreateShort(ctx, apiKey, originalUrl, expiryDate)
}

func (m *metrics) GetOriginalUrl(ctx context.Context, shortenedId string) (originalUrl string, err error) {
	defer func(begin time.Time) {
		elapsedTime := float64(time.Since(begin).Milliseconds())
		labels := prometheus.Labels{"method": "GetOriginalUrl", "error": fmt.Sprint(err != nil)}
		m.rc.With(labels).Add(1)
		m.rs.With(labels).Observe(elapsedTime)
	}(time.Now())
	return m.next.GetOriginalUrl(ctx, shortenedId)
}

func (m *metrics) Describe(ctx context.Context, shortenedId string) (resp *repository.ShortenedIdResponse, err error) {
	defer func(begin time.Time) {
		elapsedTime := float64(time.Since(begin).Milliseconds())
		labels := prometheus.Labels{"method": "Describe", "error": fmt.Sprint(err != nil)}
		m.rc.With(labels).Add(1)
		m.rs.With(labels).Observe(elapsedTime)
	}(time.Now())
	return m.next.Describe(ctx, shortenedId)
}

func (m *metrics) Delete(ctx context.Context, shortenedId string) (err error) {
	defer func(begin time.Time) {
		elapsedTime := float64(time.Since(begin).Milliseconds())
		labels := prometheus.Labels{"method": "Delete", "error": fmt.Sprint(err != nil)}
		m.rc.With(labels).Add(1)
		m.rs.With(labels).Observe(elapsedTime)
	}(time.Now())
	return m.next.Delete(ctx, shortenedId)
}

func (m *metrics) IncrementStats(ctx context.Context, shortenedId string) (err error) {
	defer func(begin time.Time) {
		elapsedTime := float64(time.Since(begin).Milliseconds())
		labels := prometheus.Labels{"method": "IncrementStats", "error": fmt.Sprint(err != nil)}
		m.rc.With(labels).Add(1)
		m.rs.With(labels).Observe(elapsedTime)
	}(time.Now())
	return m.next.IncrementStats(ctx, shortenedId)
}

func (m *metrics) Stats(ctx context.Context, shortenedId string) (resp *repository.StatsResponse, err error) {
	defer func(begin time.Time) {
		elapsedTime := float64(time.Since(begin).Milliseconds())
		labels := prometheus.Labels{"method": "Stats", "error": fmt.Sprint(err != nil)}
		m.rc.With(labels).Add(1)
		m.rs.With(labels).Observe(elapsedTime)
	}(time.Now())
	return m.next.Stats(ctx, shortenedId)
}
