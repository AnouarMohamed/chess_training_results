package service

import (
	"regexp"
	"strings"
)

const (
	minUsernameLength = 3
	maxUsernameLength = 32
	minPasswordLength = 8
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func normalizeUsername(raw string) (string, error) {
	username := strings.TrimSpace(raw)
	if len(username) < minUsernameLength || len(username) > maxUsernameLength {
		return "", ErrInvalidUsername
	}
	if !usernamePattern.MatchString(username) {
		return "", ErrInvalidUsername
	}
	return username, nil
}

func validatePassword(raw string) error {
	if len(strings.TrimSpace(raw)) < minPasswordLength {
		return ErrWeakPassword
	}
	return nil
}
