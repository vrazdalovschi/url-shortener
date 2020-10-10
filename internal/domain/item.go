package domain

import "fmt"

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
