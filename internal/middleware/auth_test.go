package middleware

import (
	"context"
	"testing"
)

func TestGetUserIDFromContext_WithUserID(t *testing.T) {
	userID := int64(123)

	// Create context with user ID using the correct key
	ctx := context.WithValue(context.Background(), UserIDKey, userID)

	// Get user ID from context
	retrievedID, ok := GetUserIDFromContext(ctx)

	if !ok {
		t.Fatal("Expected user ID to be found in context")
	}

	if retrievedID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, retrievedID)
	}
}

func TestGetUserIDFromContext_WithoutUserID(t *testing.T) {
	// Create context without user ID
	ctx := context.Background()

	// Try to get user ID from context
	_, ok := GetUserIDFromContext(ctx)

	if ok {
		t.Error("Expected no user ID in context, but found one")
	}
}

func TestGetUserIDFromContext_WithWrongType(t *testing.T) {
	// Create context with wrong type
	ctx := context.WithValue(context.Background(), UserIDKey, "not-a-number")

	// Try to get user ID from context
	_, ok := GetUserIDFromContext(ctx)

	if ok {
		t.Error("Expected no user ID due to wrong type, but found one")
	}
}

func TestGetUserIDFromContext_WithZeroValue(t *testing.T) {
	userID := int64(0)

	// Create context with zero user ID
	ctx := context.WithValue(context.Background(), UserIDKey, userID)

	// Get user ID from context
	retrievedID, ok := GetUserIDFromContext(ctx)

	if !ok {
		t.Fatal("Expected user ID to be found in context (even if zero)")
	}

	if retrievedID != 0 {
		t.Errorf("Expected user ID 0, got %d", retrievedID)
	}
}

func TestGetUserIDFromContext_WithNegativeValue(t *testing.T) {
	userID := int64(-1)

	// Create context with negative user ID
	ctx := context.WithValue(context.Background(), UserIDKey, userID)

	// Get user ID from context
	retrievedID, ok := GetUserIDFromContext(ctx)

	if !ok {
		t.Fatal("Expected user ID to be found in context")
	}

	if retrievedID != -1 {
		t.Errorf("Expected user ID -1, got %d", retrievedID)
	}
}
