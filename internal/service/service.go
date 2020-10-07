package service

import (
	"context"
	"fmt"
	"github.com/vrazdalovschi/url-shortener/internal/domain"
	"github.com/vrazdalovschi/url-shortener/internal/stackerr"
	"github.com/vrazdalovschi/url-shortener/internal/storage/postgres"
	"net/url"
)

type Service interface {
	CreateShort(ctx context.Context, apiKey, originalUrl, expiryDate string) (*domain.ShortenedIdResponse, error)
	GetOriginalUrl(ctx context.Context, shortenedId string) (originalUrl string, err error)
	Describe(ctx context.Context, shortenedId string) (*domain.ShortenedIdResponse, error)
	Delete(ctx context.Context, shortenedId string) error
	IncrementStats(ctx context.Context, shortenedId string) error
	Stats(ctx context.Context, shortenedId string) (*domain.StatsResponse, error)
}

func NewService(st postgres.Service) Service {
	return &service{st: st}
}

type service struct {
	st postgres.Service
}

func (s *service) CreateShort(ctx context.Context, apiKey, originalUrl, expiryDate string) (*domain.ShortenedIdResponse, error) {
	_, err := url.Parse(originalUrl)
	if err != nil {
		err = domain.Error{Message: fmt.Sprintf("Invalid originalUrl %v", err), ErrorCode: 400}
		return nil, stackerr.Wrap(err)
	}

	shortUrl := GenerateShortUrl(originalUrl)
	if err := s.st.Save(ctx, apiKey, originalUrl, shortUrl, expiryDate); err != nil {
		return nil, stackerr.Wrap(err)
	}
	return &domain.ShortenedIdResponse{
		ApiKey:      apiKey,
		OriginalURL: originalUrl,
		ShortenedId: shortUrl,
		ExpiryDate:  expiryDate,
	}, nil
}

func (s *service) GetOriginalUrl(ctx context.Context, shortenedId string) (originalUrl string, err error) {
	originalUrl, err = s.st.Load(ctx, shortenedId)
	if err != nil {
		return "", stackerr.Wrap(err)
	}
	return originalUrl, nil
}

func (s *service) Describe(ctx context.Context, shortenedId string) (*domain.ShortenedIdResponse, error) {
	item, err := s.st.Describe(ctx, shortenedId)
	return item, stackerr.Wrap(err)
}

func (s *service) Delete(ctx context.Context, shortenedId string) error {
	return s.st.Delete(ctx, shortenedId)
}

func (s *service) IncrementStats(ctx context.Context, shortenedId string) error {
	err := s.st.Increment(ctx, shortenedId)
	return stackerr.Wrap(err)
}

func (s *service) Stats(ctx context.Context, shortenedId string) (*domain.StatsResponse, error) {
	stats, err := s.st.Stats(ctx, shortenedId)
	return stats, stackerr.Wrap(err)
}
