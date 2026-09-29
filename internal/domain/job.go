package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

var ErrInvalidName = errors.New("name must contain 1 to 120 characters")

type Job struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func NormalizeName(value string) (string, error) {
	name := strings.TrimSpace(value)
	if count := utf8.RuneCountInString(name); count < 1 || count > 120 {
		return "", ErrInvalidName
	}
	return name, nil
}
