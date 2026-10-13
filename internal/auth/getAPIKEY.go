package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	key_params := headers.Get("Authorization")
	if key_params == "" {
		return "", errors.New("no APIkey found")
	}
	words := strings.Split(key_params, " ")
	if len(words) < 2 || words[0] != "ApiKey" {
		return "", errors.New("no APIkey found")
	}
	return words[1], nil
}
