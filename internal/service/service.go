package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/vrazdalovschi/url-shortener/internal/domain"
	"github.com/vrazdalovschi/url-shortener/internal/repository"
	"github.com/vrazdalovschi/url-shortener/internal/stackerr"
	"net/url"
	"time"
)

type Service interface {
	CreateShort(ctx context.Context, apiKey, originalUrl, expiryDate string) (*repository.ShortenedIdResponse, error)
	GetOriginalUrl(ctx context.Context, shortenedId string) (originalUrl string, err error)
	Describe(ctx context.Context, shortenedId string) (*repository.ShortenedIdResponse, error)
	Delete(ctx context.Context, shortenedId string) error
	IncrementStats(ctx context.Context, shortenedId string) error
	Stats(ctx context.Context, shortenedId string) (*repository.StatsResponse, error)
}

func NewService(st repository.Repository) Service {
	return &service{st: st}
}

type service struct {
	st repository.Repository
}

const timeStampFormat = "2006-01-02"

// Create short url for apiKey, originalUrl, expiryDate
// Mandatory parameter is originalUrl
// Default value for expiryDate is one year
func (s *service) CreateShort(ctx context.Context, apiKey, originalUrl, expiryDate string) (*repository.ShortenedIdResponse, error) {
	if !isValidUrl(originalUrl) {
		err := domain.Error{Message: fmt.Sprintf("Invalid originalUrl %s", originalUrl), ErrorCode: 400}
		return nil, stackerr.Wrap(err)
	}
	if apiKey == "" {
		apiKey = uuid.New().String()
	}
	//Ignore invalid expiryDate, set default: +1 year
	if _, err := time.Parse(timeStampFormat, expiryDate); err != nil {
		expiryDate = time.Now().AddDate(1, 0, 0).Format(timeStampFormat)
	}
	shortUrl := GenerateShortUrl(apiKey, originalUrl)
	if err := s.st.Save(ctx, apiKey, originalUrl, shortUrl, expiryDate); err != nil {
		return nil, stackerr.Wrap(err)
	}
	return &repository.ShortenedIdResponse{
		ApiKey:      apiKey,
		OriginalURL: originalUrl,
		ShortenedId: shortUrl,
		ExpiryDate:  expiryDate,
	}, nil
}

func isValidUrl(str string) bool {
	u, err := url.ParseRequestURI(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (s *service) GetOriginalUrl(ctx context.Context, shortenedId string) (originalUrl string, err error) {
	originalUrl, err = s.st.Load(ctx, shortenedId)
	return originalUrl, stackerr.Wrap(err)
}

func (s *service) Describe(ctx context.Context, shortenedId string) (*repository.ShortenedIdResponse, error) {
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

func (s *service) Stats(ctx context.Context, shortenedId string) (*repository.StatsResponse, error) {
	stats, err := s.st.Stats(ctx, shortenedId)
	return stats, stackerr.Wrap(err)
}
