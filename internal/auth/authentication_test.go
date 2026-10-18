package auth

import "testing"

func TestHashPassword(t *testing.T) {

	password := "password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned an unexpected error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned an empty token")
	}

}

func TestHashPassword_Success(t *testing.T) {
	password := "password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() returned an unexpected error: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned an empty token")
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash() returned an unexpected error: %v", err)
	}

	if !match {
		t.Fatalf("CheckPasswordHash() password did not match hash: %v", err)
	}
}

func TestHashPassword_SaltUniqueness(t *testing.T) {
	password := "RepeatPassword456!"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("first hash generation failed: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("second hash generation failed: %v", err)
	}

	if hash1 == hash2 {
		t.Errorf("expected unique hashes due to random salting, but both were identical: %s", hash1)
	}

	match1, _ := CheckPasswordHash(password, hash1)
	match2, _ := CheckPasswordHash(password, hash2)

	if !match1 || !match2 {
		t.Errorf("both unique hashes should successfully verify against the original password")
	}
}

func TestCheckPasswordHash_Match(t *testing.T) {
	password := "SecurePassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to generate helper hash: %v", err)
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !match {
		t.Errorf("expected match to be true for the correct password")
	}
}

func TestCheckPasswordHash_NoMatch(t *testing.T) {
	password := "SecurePassword123!"
	wrongPassword := "WrongPassword999!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to generate helper hash: %v", err)
	}

	match, err := CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("expected no error even for a password mismatch, got %v", err)
	}

	if match {
		t.Errorf("expected match to be false for an incorrect password")
	}
}

func TestCheckPasswordHash_InvalidHashError(t *testing.T) {
	password := "SecurePassword123!"
	invalidHash := "this-is-not-a-valid-argon2id-hash-string"

	match, err := CheckPasswordHash(password, invalidHash)
	if err == nil {
		t.Errorf("expected an error due to malformed hash format, got nil")
	}

	if match {
		t.Errorf("expected match to be false when an error occurs")
	}
}
