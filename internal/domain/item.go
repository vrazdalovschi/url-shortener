package domain

type Item struct {
	ApiKey      string `json:"apiKey"`
	OriginalURL string `json:"originalUrl"`
	ShortenedId string `json:"shortenedId"`
	ExpiryDate  string `json:"expiryDate"`
}
