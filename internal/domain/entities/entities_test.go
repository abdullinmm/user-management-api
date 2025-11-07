package entities

import (
	"testing"
	"time"
)

func TestUser_Creation(t *testing.T) {
	user := &User{
		ID:        1,
		Username:  "testuser",
		CreatedAt: time.Now(),
	}
	
	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}
	
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}
	
	if user.ReferrerID != nil {
		t.Error("Expected nil ReferrerID")
	}
}

func TestUser_WithReferrer(t *testing.T) {
	referrerID := int64(10)
	user := &User{
		ID:         2,
		Username:   "referred_user",
		ReferrerID: &referrerID,
		CreatedAt:  time.Now(),
	}
	
	if user.ReferrerID == nil {
		t.Fatal("Expected ReferrerID, got nil")
	}
	
	if *user.ReferrerID != 10 {
		t.Errorf("Expected ReferrerID 10, got %d", *user.ReferrerID)
	}
}

func TestTask_Validation(t *testing.T) {
	task := &Task{
		ID:           1,
		Code:         "TASK_REGISTER",
		Title:        "Complete Registration",
		RewardPoints: 100,
		IsActive:     true,
	}
	
	if task.RewardPoints != 100 {
		t.Errorf("Expected reward 100, got %d", task.RewardPoints)
	}
	
	if !task.IsActive {
		t.Error("Expected task to be active")
	}
}

func TestBalance_Operations(t *testing.T) {
	balance := &Balance{
		UserID: 1,
		Points: 0,
	}
	
	if balance.Points != 0 {
		t.Errorf("Expected 0 points, got %d", balance.Points)
	}
	
	// Simulate adding points
	balance.Points += 100
	
	if balance.Points != 100 {
		t.Errorf("Expected 100 points after addition, got %d", balance.Points)
	}
}

func TestTransaction_CreationWithTask(t *testing.T) {
	refID := int64(5)
	refType := "task"
	
	tx := &Transaction{
		ID:            1,
		UserID:        1,
		Delta:         50,
		Reason:        "task_completion",
		ReferenceID:   &refID,
		ReferenceType: &refType,
		CreatedAt:     time.Now(),
	}
	
	if tx.Delta != 50 {
		t.Errorf("Expected delta 50, got %d", tx.Delta)
	}
	
	if tx.ReferenceID == nil {
		t.Fatal("Expected ReferenceID, got nil")
	}
	
	if *tx.ReferenceID != 5 {
		t.Errorf("Expected ReferenceID 5, got %d", *tx.ReferenceID)
	}
	
	if tx.Reason != "task_completion" {
		t.Errorf("Expected reason 'task_completion', got %s", tx.Reason)
	}
}

func TestTransaction_NegativeDelta(t *testing.T) {
	tx := &Transaction{
		ID:        2,
		UserID:    1,
		Delta:     -25,
		Reason:    "penalty",
		CreatedAt: time.Now(),
	}
	
	if tx.Delta != -25 {
		t.Errorf("Expected delta -25, got %d", tx.Delta)
	}
	
	if tx.ReferenceID != nil {
		t.Error("Expected nil ReferenceID")
	}
}
