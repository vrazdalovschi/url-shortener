package storage

import (
	"context"
	"github.com/vrazdalovschi/url-shortener/internal/storage/postgres"
)

type Service interface {
	Save(ctx context.Context, apiKey, originalUrl, shortenedId, expiryDate string) error
	Load(ctx context.Context, shortenedId string) (originalUrl string, err error)
	Describe(ctx context.Context, shortenedId string) (*Item, error)
	Close() error
}

type Configuration struct {
	Host, Port, User, Password, DbName string
}

func GetActive(cfg Configuration) (Service, error) {
	return postgres.New(cfg)
}

type Item struct {
	ApiKey      string `json:"apiKey"`
	OriginalURL string `json:"originalUrl"`
	ShortenedId string `json:"shortenedId"`
	Enabled     bool   `json:"enabled"`
	ExpiryDate  string `json:"expiryDate"`
}
