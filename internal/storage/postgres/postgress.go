package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vrazdalovschi/url-shortener/internal/domain"
	"github.com/vrazdalovschi/url-shortener/internal/stackerr"
	"time"

	// This loads the postgres drivers.
	_ "github.com/lib/pq"
)

type Service interface {
	Save(ctx context.Context, apiKey, originalUrl, shortenedId, expiryDate string) error
	Load(ctx context.Context, shortenedId string) (originalUrl string, err error)
	Describe(ctx context.Context, shortenedId string) (*domain.Item, error)
	Delete(ctx context.Context, shortenedId string) error
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
		return nil, stackerr.Wrap(err)
	}

	// Ping to connection
	err = db.Ping()
	if err != nil {
		return nil, stackerr.Wrap(err)
	}

	// Create table if not exists
	strQuery := "CREATE TABLE IF NOT EXISTS url (shortenedId VARCHAR NOT NULL UNIQUE, originalUrl VARCHAR not NULL, " +
		"apiKey VARCHAR not NULL, creationTime timestamp not NULL, expirationDate timestamp not NULL);"

	if _, err = db.Exec(strQuery); err != nil {
		return nil, stackerr.Wrap(err)
	}
	return &postgres{db}, nil
}

type postgres struct{ db *sql.DB }

const timeStampFormat = "2006-01-02"

func (p *postgres) Save(ctx context.Context, apiKey, originalUrl, shortenedId, expiryDate string) error {
	if _, err := time.Parse(timeStampFormat, expiryDate); err != nil {
		expiryDate = time.Now().AddDate(1, 0, 0).Format(timeStampFormat)
	}
	_, err := p.db.ExecContext(ctx, "INSERT INTO url (shortenedId, originalUrl, apiKey, creationTime, expirationDate) VALUES ($1, $2, $3, NOW(), TO_TIMESTAMP($4, 'YYYY-MM-DD'))",
		shortenedId, originalUrl, apiKey, expiryDate)
	return stackerr.Wrap(err)
}

func (p *postgres) Delete(ctx context.Context, shortenedId string) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM url WHERE shortenedId = $1", shortenedId)
	return stackerr.Wrap(err)
}

func (p *postgres) Load(ctx context.Context, shortenedId string) (originalUrl string, err error) {
	err = p.db.QueryRowContext(ctx, "SELECT originalUrl FROM url WHERE shortenedId = $1", shortenedId).Scan(&originalUrl)
	if err != nil {
		return "", stackerr.Wrap(err)
	}
	return originalUrl, nil
}

func (p *postgres) Describe(ctx context.Context, shortenedId string) (*domain.Item, error) {
	item := domain.Item{ShortenedId: shortenedId}
	query := "SELECT originalUrl, apiKey, expirationDate FROM url WHERE shortenedId = $1"
	err := p.db.QueryRowContext(ctx, query, shortenedId).Scan(&item.OriginalURL, &item.ApiKey, &item.ExpiryDate)
	if err != nil {
		return nil, stackerr.Wrap(err)
	}
	return &item, nil
}

func (p *postgres) Close() error { return p.db.Close() }
