package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {

	userID := uuid.New()
	tokenSecret := "9b1deb4d3b7d4f82af647d6928e1d51c70e28f32cf4b46c6bbf4479e0094e9f3"
	expiresIn := 5 * time.Second
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)

	if err != nil {
		t.Fatalf("MakeJWT() returned an unexpected error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("MakeJWT() returned an empty token")
	}

}

func TestValidateJWT(t *testing.T) {

	userID := uuid.New()
	tokenSecret := "9b1deb4d3b7d4f82af647d6928e1d51c70e28f32cf4b46c6bbf4479e0094e9f3"
	expiresIn := 5 * time.Second

	tokenString, err := MakeJWT(uuid.UUID(userID), tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() returned an unexpected error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("MakeJWT() returned an empty token")
	}

	gotUserID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT() returned an unexpected error: %v", err)
	}
	if gotUserID != userID {

		t.Errorf("ValidateJWT() userID = %v, want %v", gotUserID, userID)
	}

}

func TestValidateJWTExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "9b1deb4d3b7d4f82af647d6928e1d51c70e28f32cf4b46c6bbf4479e0094e9f3"
	expiresIn := -5 * time.Second

	tokenString, err := MakeJWT(uuid.UUID(userID), tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() returned an unexpected error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("MakeJWT() returned an empty token")
	}

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatal("ValidateJWT() expected an error for an expired token, got nil")
	}

}

func TestValidateJWTWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "9b1deb4d3b7d4f82af647d6928e1d51c70e28f32cf4b46c6bbf4479e0094e9f3"
	wrongSecret := "8b1deb4d3b7d4f82af647d6928e1d51c70e28f32cf4b46c6bbf4479e0094e9f2"
	expiresIn := 5 * time.Second

	tokenString, err := MakeJWT(uuid.UUID(userID), tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() returned an unexpected error: %v", err)
	}
	if tokenString == "" {
		t.Fatal("MakeJWT() returned an empty token")
	}

	_, err = ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Fatal("ValidateJWT() expected an error for wrong secret key, got nil")
	}
}

func TestGetBearerTokenExists(t *testing.T) {

	header := http.Header{}

	header.Set("Authorization", "Bearer FNisyAvh914a")

	_, err := GetBearerToken(header)
	if err != nil {
		t.Fatal("GetBearerToken() expected no error, got an error:", err)
	}

}

func TestBearerTokenDoesNotExist(t *testing.T) {
	header := http.Header{}

	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatal("GetBearerToken() expected an error for empty token ,got nil:", err)
	}
}
