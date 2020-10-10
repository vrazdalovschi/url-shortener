package repository

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, apiKey, originalUrl, shortenedId, expiryDate string) error
	Load(ctx context.Context, shortenedId string) (originalUrl string, err error)
	Describe(ctx context.Context, shortenedId string) (*ShortenedIdResponse, error)
	Delete(ctx context.Context, shortenedId string) error
	Increment(ctx context.Context, shortenedId string) error
	Stats(ctx context.Context, shortenedId string) (*StatsResponse, error)
	Close() error
}

type ShortenedIdResponse struct {
	ApiKey      string `json:"apiKey"`
	OriginalURL string `json:"originalUrl"`
	ShortenedId string `json:"shortenedId"`
	ExpiryDate  string `json:"expiryDate"`
}

type StatsResponse struct {
	ShortenedId  string `json:"shortenedId"`
	Redirects    int    `json:"redirects"`
	LastRedirect string `json:"lastRedirect"`
}

type Configuration struct {
	Host, Port, User, Password, DbName string
}
