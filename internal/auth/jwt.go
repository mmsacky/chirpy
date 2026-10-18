package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	now := time.Now().UTC()

	registeredClaims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}

	keyFunc := func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		return uuid.Nil, err
	}

	stringUserID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(stringUserID)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil

}

func GetBearerToken(headers http.Header) (string, error) {

	bearerToken := headers.Get("Authorization")

	if bearerToken == "" {
		return "", errors.New("no authorization header included")
	}

	tokenString, found := strings.CutPrefix(bearerToken, "Bearer ")
	if !found {
		return "", errors.New("malformed authorization header")
	}
	return strings.TrimSpace(tokenString), nil

}
