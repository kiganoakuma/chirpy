package auth

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "secret_password"

	// Test successful hashing
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Verify the hash is not empty
	if hash == "" {
		t.Error("Expected non-empty hash, got empty string")
	}

	// Verify the hash is different from the original password
	if hash == password {
		t.Error("Hash should not be the same as the original password")
	}

	// Verify the hash is a valid bcrypt hash
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		t.Errorf("Generated hash is not valid: %v", err)
	}
}

func TestCheckPassword(t *testing.T) {
	password := "test_password"
	wrongPassword := "wrong_password"

	// First, create a hash
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password for test: %v", err)
	}

	// Test correct password
	err = CheckPassword(hash, password)
	if err != nil {
		t.Errorf("CheckPassword failed with correct password: %v", err)
	}

	// Test incorrect password
	err = CheckPassword(hash, wrongPassword)
	if err == nil {
		t.Error("CheckPassword should fail with incorrect password")
	}
}
