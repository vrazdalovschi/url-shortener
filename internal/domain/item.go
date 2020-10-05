package domain

type Item struct {
	ApiKey      string `json:"apiKey"`
	OriginalURL string `json:"originalUrl"`
	ShortenedId string `json:"shortenedId"`
	Enabled     bool   `json:"enabled"`
	ExpiryDate  string `json:"expiryDate"`
}
