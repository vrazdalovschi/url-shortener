package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vrazdalovschi/url-shortener/internal/mocks"
	"github.com/vrazdalovschi/url-shortener/internal/repository"
	"testing"
	"time"
)

func setup() (svc Service, mockedDb *mocks.MockedPostgres) {
	mockedDb = mocks.NewMockedPostgres()
	svc = NewService(mockedDb)
	return svc, mockedDb
}

func TestService_CreateShort_WrongUrls(t *testing.T) {
	svc, db := setup()
	db.On("Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	testCases := []struct {
		URL         string
		expectError bool
	}{
		{"google.com", true},
		{"http//google.com", true},
		{"/foo/bar", true},
		{"http://google.com", false},
	}
	for _, testCase := range testCases {
		_, err := svc.CreateShort(context.Background(), "", testCase.URL, "")
		if testCase.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestService_CreateShort(t *testing.T) {
	svc, db := setup()
	originalUrl := "http://google.com"
	var (
		expectedApikey   string
		expectedShortKey string
	)
	argumentMatcherApiKey := mock.MatchedBy(func(apiKey string) bool {
		expectedApikey = apiKey
		return assert.NotEmpty(t, apiKey)
	})
	argumentMatcherShortKey := mock.MatchedBy(func(shortKey string) bool {
		expectedShortKey = shortKey
		return assert.NotEmpty(t, shortKey)
	})
	oneYearPlus := time.Now().AddDate(1, 0, 0).Format(timeStampFormat)
	db.On("Save", mock.Anything, argumentMatcherApiKey, originalUrl, argumentMatcherShortKey, oneYearPlus).Return(nil)

	actualRes, err := svc.CreateShort(context.Background(), "", originalUrl, "")
	require.NoError(t, err)

	expectedRes := &repository.ShortenedIdResponse{
		ApiKey:      expectedApikey,
		OriginalURL: originalUrl,
		ShortenedId: expectedShortKey,
		ExpiryDate:  oneYearPlus,
	}
	require.Equal(t, expectedRes, actualRes)
}
