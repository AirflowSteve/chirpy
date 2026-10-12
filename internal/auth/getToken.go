package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth_param := headers.Get("Authorization")
	if auth_param == "" {
		return "", errors.New("no autherization header found")
	}
	params := strings.Split(auth_param, " ")
	if len(params) != 2 {
		return "", errors.New("couldn't get the auth_token")
	}
	bearer := params[1]
	return bearer, nil
}
