package domain

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

type StatsResponse struct {
	ShortenedId  string `json:"shortenedId"`
	Redirects    int    `json:"redirects"`
	LastRedirect string `json:"lastRedirect"`
}
