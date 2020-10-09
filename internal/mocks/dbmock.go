package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/vrazdalovschi/url-shortener/internal/domain"
)

type MockedPostgres struct {
	mock.Mock
}

func NewMockedPostgres() *MockedPostgres {
	return &MockedPostgres{}
}

func (m *MockedPostgres) Save(ctx context.Context, apiKey, originalUrl, shortenedId, expiryDate string) error {
	called := m.Called(ctx, apiKey, originalUrl, shortenedId, expiryDate)
	return called.Error(0)
}

func (m *MockedPostgres) Load(ctx context.Context, shortenedId string) (originalUrl string, err error) {
	called := m.Called(ctx, originalUrl, shortenedId)
	return called.String(0), called.Error(1)
}

func (m *MockedPostgres) Describe(ctx context.Context, shortenedId string) (*domain.ShortenedIdResponse, error) {
	called := m.Called(ctx, shortenedId)
	return called.Get(0).(*domain.ShortenedIdResponse), called.Error(1)
}

func (m *MockedPostgres) Delete(ctx context.Context, shortenedId string) error {
	called := m.Called(ctx, shortenedId)
	return called.Error(0)
}

func (m *MockedPostgres) Increment(ctx context.Context, shortenedId string) error {
	called := m.Called(ctx, shortenedId)
	return called.Error(0)
}
func (m *MockedPostgres) Stats(ctx context.Context, shortenedId string) (*domain.StatsResponse, error) {
	called := m.Called(ctx, shortenedId)
	return called.Get(0).(*domain.StatsResponse), called.Error(1)
}

func (m *MockedPostgres) Close() error {
	called := m.Called()
	return called.Error(0)
}
