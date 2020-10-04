package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/vrazdalovschi/url-shortener/external/github.com/lytics/base62"
)

/*
Encode apiKey with OriginalUrl together to not create the same short url for different users
*/
func GenerateShortUrl(apiKey, originalUrl string) string {
	if apiKey == "" {
		apiKey = uuid.New().String()
	}
	compoundKey := fmt.Sprintf("%s_%s", apiKey, originalUrl)
	generatedKey := hashToMd5(compoundKey)[0:7]
	encodedResult := base62.StdEncoding.EncodeToString([]byte(generatedKey))
	return encodedResult
}

func hashToMd5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
