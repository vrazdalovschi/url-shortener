package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/vrazdalovschi/url-shortener/external/github.com/lytics/base62"
	"time"
)

// Generate random shortId.
// Should be unique each time and shouldn't repeat for the same url
// Limitations: the same parameters and in the same time, then would be generating the same shortId.
func GenerateShortUrl(apiKey, originalUrl string) string {
	compoundKey := fmt.Sprintf("%s_%s_%s", apiKey, originalUrl, time.Now().String())
	generatedKey := hash(compoundKey)[0:6]
	encodedResult := base62.StdEncoding.EncodeToString([]byte(generatedKey))
	return encodedResult
}

func hash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
