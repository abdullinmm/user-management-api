package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTTokenCreation(t *testing.T) {
	secret := []byte("test-secret-key")
	userID := int64(123)

	// Create token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)

	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	if tokenString == "" {
		t.Fatal("Token string is empty")
	}
}

func TestJWTTokenParsing(t *testing.T) {
	secret := []byte("test-secret-key")
	userID := int64(456)

	// Create token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secret)

	// Parse token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if !parsedToken.Valid {
		t.Fatal("Token is invalid")
	}

	// Check claims
	parsedClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to get claims")
	}

	if parsedClaims["user_id"].(float64) != float64(userID) {
		t.Errorf("Expected user_id %d, got %v", userID, parsedClaims["user_id"])
	}
}

func TestJWTExpiredToken(t *testing.T) {
	secret := []byte("test-secret-key")

	// Create expired token
	claims := jwt.MapClaims{
		"user_id": int64(789),
		"exp":     time.Now().Add(-1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(secret)

	// Try to parse expired token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	// Should fail or be invalid
	if err == nil && parsedToken.Valid {
		t.Error("Expected error or invalid token for expired token")
	}
}
