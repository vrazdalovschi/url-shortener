package serivcemiddleware

import (
	"context"
	"github.com/vrazdalovschi/url-shortener/internal/domain"
	"github.com/vrazdalovschi/url-shortener/internal/service"
	"log"
	"time"
)

type loggingMiddleware struct {
	next service.Service
}

func NewLogging() Middleware {
	return func(s service.Service) service.Service {
		return &loggingMiddleware{
			next: s,
		}
	}
}

func (l *loggingMiddleware) CreateShort(ctx context.Context, apiKey, originalUrl, expiryDate string) (resp *domain.ShortenedIdResponse, err error) {
	defer func(begin time.Time) {
		log.Println("method", "CreateShort",
			"apiKey", apiKey,
			"originalUrl", originalUrl,
			"expiryDate", expiryDate,
			"err", err,
			"resp", resp,
			"duration", time.Since(begin))
	}(time.Now())

	return l.next.CreateShort(ctx, apiKey, originalUrl, expiryDate)
}

func (l *loggingMiddleware) GetOriginalUrl(ctx context.Context, shortenedId string) (originalUrl string, err error) {
	defer func(begin time.Time) {
		log.Println("method", "GetOriginalUrl",
			"shortenedId", shortenedId,
			"err", err,
			"originalUrl", originalUrl,
			"duration", time.Since(begin))
	}(time.Now())

	return l.next.GetOriginalUrl(ctx, shortenedId)
}

func (l *loggingMiddleware) Describe(ctx context.Context, shortenedId string) (resp *domain.ShortenedIdResponse, err error) {
	defer func(begin time.Time) {
		log.Println("method", "Describe",
			"shortenedId", shortenedId,
			"err", err,
			"resp", resp,
			"duration", time.Since(begin))
	}(time.Now())

	return l.next.Describe(ctx, shortenedId)
}

func (l *loggingMiddleware) Delete(ctx context.Context, shortenedId string) (err error) {
	defer func(begin time.Time) {
		log.Println("method", "Delete",
			"shortenedId", shortenedId,
			"err", err,
			"duration", time.Since(begin))
	}(time.Now())

	return l.next.Delete(ctx, shortenedId)
}
