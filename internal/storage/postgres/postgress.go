package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vrazdalovschi/url-shortener/internal/domain"
	"time"

	// This loads the postgres drivers.
	_ "github.com/lib/pq"
)

type Service interface {
	Save(ctx context.Context, apiKey, originalUrl, shortenedId, expiryDate string) error
	Load(ctx context.Context, shortenedId string) (originalUrl string, err error)
	Describe(ctx context.Context, shortenedId string) (*domain.Item, error)
	Close() error
}

type Configuration struct {
	Host, Port, User, Password, DbName string
}

// New returns a postgres backed storage service.
func New(cfg Configuration) (Service, error) {
	// Connect postgres
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)

	db, err := sql.Open("postgres", connect)
	if err != nil {
		return nil, err
	}

	// Ping to connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	strQuery := "CREATE TABLE IF NOT EXISTS url (shortenedId VARCHAR NOT NULL UNIQUE, originalUrl VARCHAR not NULL, " +
		"apiKey VARCHAR not NULL, creationTime datetime not NULL, expirationDate datetime not NULL);"

	_, err = db.Exec(strQuery)
	if err != nil {
		return nil, err
	}
	return &postgres{db}, nil
}

type postgres struct{ db *sql.DB }

func (p *postgres) Save(ctx context.Context, apiKey, originalUrl, shortenedId, expiryDate string) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO url (shortenedId, originalUrl, apiKey, creationTime, expirationTime) VALUES ($1, $2, $3, $4, $5)",
		shortenedId, originalUrl, apiKey, time.Now().String(), expiryDate)
	return err
}

func (p *postgres) Load(ctx context.Context, shortenedId string) (originalUrl string, err error) {
	res, err := p.db.QueryContext(ctx, "SELECT originalUrl FROM url WHERE shortenedId = $1 limit 1", shortenedId)
	if err != nil {
		return "", err
	}
	item := domain.Item{ShortenedId: shortenedId}
	if err = res.Scan(&item.OriginalURL); res != nil {
		return "", err
	}
	return item.OriginalURL, nil
}

func (p *postgres) Describe(ctx context.Context, shortenedId string) (*domain.Item, error) {
	res, err := p.db.QueryContext(ctx, "SELECT originalUrl, apiKey, enable, expiryDate FROM url WHERE shortenedId = $1 limit 1", shortenedId)
	if err != nil {
		return nil, err
	}
	item := domain.Item{ShortenedId: shortenedId}
	if err = res.Scan(&item.OriginalURL, &item.ApiKey, &item.Enabled, &item.ExpiryDate); res != nil {
		return nil, err
	}
	return &item, nil
}

func (p *postgres) Close() error { return p.db.Close() }
