package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if token, hasPrefix := strings.CutPrefix(authorization, "Bearer "); hasPrefix && len(token) > 0 {
		return token, nil
	}
	return "", errors.New("Could not parse authorization header")
}
