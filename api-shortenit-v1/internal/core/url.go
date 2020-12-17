package core

// ShortenURLRequest ...
type ShortenURLRequest struct {
	OriginalURL string `json:"originUrl,omitempty"`
	CustomAlias string `json:"customAlias,omitempty"`
	UserEmail   string `json:"userEmail,omitempty"`
}

// ShortenURLResponse ...
type ShortenURLResponse struct {
	URL string `json:"url,omitempty"`
}
