package service

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGenerateShortUrl(t *testing.T) {
	apiKey := uuid.New().String()
	testUrl := "http://www.google.com/findme?key=qqq&value=zzzz"
	var empty interface{}
	set := map[string]interface{}{}
	testNumber := 1000
	for i := 0; i < testNumber; i++ {
		generated := GenerateShortUrl(apiKey, testUrl)
		set[generated] = empty
		time.Sleep(time.Nanosecond * 1)
	}
	require.Equal(t, testNumber, len(set))
}
