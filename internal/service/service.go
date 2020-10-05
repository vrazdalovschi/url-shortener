package service

import (
	"context"
	"github.com/vrazdalovschi/url-shortener/internal/domain"
	"github.com/vrazdalovschi/url-shortener/internal/stackerr"
	"github.com/vrazdalovschi/url-shortener/internal/storage/postgres"
)

type Service interface {
	CreateShort(ctx context.Context, apiKey, originalUrl, expiryDate string) (string, error)
	GetOriginalUrl(ctx context.Context, shortenedId string) (originalUrl string, err error)
	Describe(ctx context.Context, shortenedId string) (*domain.Item, error)
	Delete(ctx context.Context, shortenedId string) error
}

func NewService(st postgres.Service) Service {
	return &service{st: st}
}

type service struct {
	st postgres.Service
}

func (s *service) CreateShort(ctx context.Context, apiKey, originalUrl, expiryDate string) (string, error) {
	shortUrl := GenerateShortUrl(originalUrl)
	if err := s.st.Save(ctx, apiKey, originalUrl, shortUrl, expiryDate); err != nil {
		return "", stackerr.Wrap(err)
	}
	return shortUrl, nil
}

func (s *service) GetOriginalUrl(ctx context.Context, shortenedId string) (originalUrl string, err error) {
	url, err := s.st.Load(ctx, shortenedId)
	if err != nil {
		return "", stackerr.Wrap(err)
	}
	return url, nil
}

func (s *service) Describe(ctx context.Context, shortenedId string) (*domain.Item, error) {
	item, err := s.st.Describe(ctx, shortenedId)
	if err != nil {
		return nil, stackerr.Wrap(err)
	}
	return item, nil
}

func (s *service) Delete(ctx context.Context, shortenedId string) error {
	return s.st.Delete(ctx, shortenedId)
}
