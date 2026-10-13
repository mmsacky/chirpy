package auth

import (
	"testing"
)

// TestHashPassword_Success verifies that a password is successfully hashed,
// returns a non-empty string, and can be verified back using argon2id.
func TestHashPassword_Success(t *testing.T) {
	password := "SecurePassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatalf("expected a non-empty hash string, got an empty string")
	}

	// Verify that the generated hash correctly matches the original password
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("unexpected error when comparing password and hash: %v", err)
	}

	if !match {
		t.Errorf("expected the generated hash to match the plain-text password")
	}
}

// TestHashPassword_SaltUniqueness ensures that hashing the same password
// multiple times produces different hash strings due to cryptographically-secure random salts.
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

	// Both distinct hashes should still successfully verify against the password
	match1, _ := CheckPasswordHash(password, hash1)
	match2, _ := CheckPasswordHash(password, hash2)

	if !match1 || !match2 {
		t.Errorf("both unique hashes should successfully verify against the original password")
	}
}

// TestCheckPasswordHash_Match verifies that a correct password returns true and no error.
func TestCheckPasswordHash_Match(t *testing.T) {
	password := "SecurePassword123!"

	// Generate a valid hash for testing purposes
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

// TestCheckPasswordHash_NoMatch verifies that an incorrect password returns false and no error.
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

// TestCheckPasswordHash_InvalidHashError verifies that a malformed hash string
// correctly triggers and returns an error.
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
