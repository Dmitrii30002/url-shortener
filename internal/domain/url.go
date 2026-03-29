package domain

import "time"

type URL struct {
	ID          int64      `json:"id"`
	OriginalURL string     `json:"original_url"`
	ShortURL    string     `json:"short_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
