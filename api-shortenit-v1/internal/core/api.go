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

type URLResponse ShortenURLResponse

func (r *ShortenURLRequest) Size() int64 {
	return int64(len(r.CustomAlias) + len(r.OriginalURL) + len(r.UserEmail))
}