package service

import (
	"context"
	"github.com/vrazdalovschi/url-shortener/internal/domain"
	"github.com/vrazdalovschi/url-shortener/internal/storage/postgres"
)

type Service interface {
	CreateShort(ctx context.Context, apiKey, originalUrl, expiryDate string) error
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

func (s *service) CreateShort(ctx context.Context, apiKey, originalUrl, expiryDate string) error {
	panic("implement me")
}

func (s *service) GetOriginalUrl(ctx context.Context, shortenedId string) (originalUrl string, err error) {
	panic("implement me")
}

func (s *service) Describe(ctx context.Context, shortenedId string) (*domain.Item, error) {
	panic("implement me")
}

func (s *service) Delete(ctx context.Context, shortenedId string) error {
	panic("implement me")
}
