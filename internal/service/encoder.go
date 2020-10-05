package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/vrazdalovschi/url-shortener/external/github.com/lytics/base62"
	"time"
)

func GenerateShortUrl(originalUrl string) string {
	compoundKey := fmt.Sprintf("%s_%s", originalUrl, time.Now().String())
	generatedKey := hash(compoundKey)[0:6]
	encodedResult := base62.StdEncoding.EncodeToString([]byte(generatedKey))
	return encodedResult
}

func hash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
