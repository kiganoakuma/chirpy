package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret_key"
	expiresIn := time.Hour * 24

	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	returnedID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	if returnedID != userID {
		t.Errorf("UserID mismatch. Expected: %s, Got: %s", userID, returnedID)
	}
}

func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"
	expiration := time.Millisecond * 10

	tokenString, err := MakeJWT(userID, secret, expiration)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	time.Sleep(time.Millisecond * 20)

	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Error("Expected error for expired token, but got nil")
	}
}

func TestInvalidSignature(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret-key"
	wrongSecret := "wrong-secret"
	expiration := time.Hour * 24

	tokenString, err := MakeJWT(userID, secret, expiration)
	if err != nil {
		t.Fatalf("Failed to  create JWT: %v", err)
	}

	_, err = ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Error("Expected error for invalid signature, but got nil")
	}
}

func TestMalformedToken(t *testing.T) {
	malformedToken := "not.a.valid.token"
	secret := "test-secret"

	_, err := ValidateJWT(malformedToken, secret)
	if err == nil {
		t.Error("Expected error for malformed token, but got nil")
	}
}

func TestTokenWithInvalidUUID(t *testing.T) {
	claims := jwt.RegisteredClaims{
		Issuer:    "Chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject:   "not-a-valid-uuid", // Invalid UUID format
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := "test-secret"

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Failed to sign test token: %v", err)
	}

	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Error("Expected error for token with invalid UUID, but got nil")
	}

}
