package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func GetAPIKey(headers http.Header) (string, error) {

	apiKey := headers.Get("Authorization")

	if apiKey == "" {
		return "", errors.New("no authorization header included")
	}

	apiKeyString, found := strings.CutPrefix(apiKey, "ApiKey ")
	if !found {
		return "", errors.New("malformed authorization header")
	}
	return strings.TrimSpace(apiKeyString), nil

}
