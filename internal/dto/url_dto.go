package dto

type CreateURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type CreateURLResponse struct {
	ShortURL string `json:"short_url"`
}

type GetOriginalURLResponse struct {
	OriginalURL string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
