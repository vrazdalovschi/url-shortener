package domain

import "fmt"

type ShortenedIdResponse struct {
	ApiKey      string `json:"apiKey"`
	OriginalURL string `json:"originalUrl"`
	ShortenedId string `json:"shortenedId"`
	ExpiryDate  string `json:"expiryDate"`
}

type CreateShortId struct {
	ApiKey      string `json:"apiKey"`
	OriginalURL string `json:"originalUrl"`
	ExpiryDate  string `json:"expiryDate"`
}

type Error struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf("Message: %s, ErrorCode: %d", e.Message, e.ErrorCode)
}

type StatsResponse struct {
	ShortenedId  string `json:"shortenedId"`
	Redirects    int    `json:"redirects"`
	LastRedirect string `json:"lastRedirect"`
}
