package domain

import "errors"

var (
	ErrNotFound    = errors.New("resource not found")
	ErrInvalidLang = errors.New("invalid language parameter")
)
