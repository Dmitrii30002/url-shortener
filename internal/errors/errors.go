package errors

import "errors"

var (
	ErrNotFound          = errors.New("url not found")
	ErrDuplicateURL      = errors.New("url already exists")
	ErrDuplicateShortURL = errors.New("short url already exists")
)
