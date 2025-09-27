package repository

import (
	"testing"
	"time"

	"share-screen/pkg/domain/entities"
)

func TestMemorySessionRepository_CreateSession(t *testing.T) {
	repo := NewMemorySessionRepository().(*MemorySessionRepository)

	expiryDuration := 30 * time.Minute
	session, err := repo.CreateSession(expiryDuration)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if session == nil {
		t.Fatal("Expected session but got nil")
	}

	if session.Token == "" {
		t.Error("Expected token but got empty string")
	}

	if session.Status != entities.SessionStatusPending {
		t.Errorf("Expected status %v but got %v", entities.SessionStatusPending, session.Status)
	}

	if time.Until(session.ExpiresAt) < 29*time.Minute {
		t.Error("Session expiry time is too short")
	}

	if time.Until(session.ExpiresAt) > 31*time.Minute {
		t.Error("Session expiry time is too long")
	}
}

func TestMemorySessionRepository_GetSession(t *testing.T) {
	repo := NewMemorySessionRepository().(*MemorySessionRepository)

	// Test getting non-existent session
	_, err := repo.GetSession("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent session")
	}

	// Create a session first
	original, err := repo.CreateSession(30 * time.Minute)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Test getting existing session
	retrieved, err := repo.GetSession(original.Token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected session but got nil")
	}

	if retrieved.Token != original.Token {
		t.Errorf("Expected token %q but got %q", original.Token, retrieved.Token)
	}

	if retrieved.Status != original.Status {
		t.Errorf("Expected status %v but got %v", original.Status, retrieved.Status)
	}

	// Test that returned session is a copy (modifications don't affect original)
	retrieved.Status = entities.SessionStatusActive
	retrievedAgain, _ := repo.GetSession(original.Token)
	if retrievedAgain.Status == entities.SessionStatusActive {
		t.Error("Session should be a copy, modifications should not persist")
	}
}

func TestMemorySessionRepository_UpdateSession(t *testing.T) {
	repo := NewMemorySessionRepository().(*MemorySessionRepository)

	// Test updating non-existent session
	nonExistentSession := &entities.Session{
		Token:     "non-existent",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Status:    entities.SessionStatusPending,
	}
	err := repo.UpdateSession(nonExistentSession)
	if err == nil {
		t.Error("Expected error for non-existent session")
	}

	// Create a session first
	original, err := repo.CreateSession(30 * time.Minute)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Update the session
	original.Status = entities.SessionStatusActive
	original.Offer = &entities.WebRTCOffer{
		Type: "offer",
		SDP:  "test-sdp",
	}

	err = repo.UpdateSession(original)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the update
	updated, err := repo.GetSession(original.Token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if updated.Status != entities.SessionStatusActive {
		t.Errorf("Expected status %v but got %v", entities.SessionStatusActive, updated.Status)
	}

	if updated.Offer == nil {
		t.Error("Expected offer but got nil")
	}

	if updated.Offer != nil && updated.Offer.Type != "offer" {
		t.Errorf("Expected offer type %q but got %q", "offer", updated.Offer.Type)
	}
}

func TestMemorySessionRepository_DeleteSession(t *testing.T) {
	repo := NewMemorySessionRepository().(*MemorySessionRepository)

	// Create a session first
	session, err := repo.CreateSession(30 * time.Minute)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Verify session exists
	_, err = repo.GetSession(session.Token)
	if err != nil {
		t.Errorf("Session should exist: %v", err)
	}

	// Delete the session
	err = repo.DeleteSession(session.Token)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify session is deleted
	_, err = repo.GetSession(session.Token)
	if err == nil {
		t.Error("Expected error for deleted session")
	}

	// Test deleting non-existent session (should not error)
	err = repo.DeleteSession("non-existent")
	if err != nil {
		t.Errorf("Unexpected error deleting non-existent session: %v", err)
	}
}

func TestMemorySessionRepository_CleanupExpiredSessions(t *testing.T) {
	repo := NewMemorySessionRepository().(*MemorySessionRepository)

	// Create some sessions with different expiry times
	now := time.Now()

	// Expired session
	expiredSession := &entities.Session{
		Token:     "expired",
		CreatedAt: now.Add(-60 * time.Minute),
		ExpiresAt: now.Add(-30 * time.Minute),
		Status:    entities.SessionStatusPending,
	}
	repo.sessions[expiredSession.Token] = expiredSession

	// Valid session
	validSession := &entities.Session{
		Token:     "valid",
		CreatedAt: now.Add(-10 * time.Minute),
		ExpiresAt: now.Add(20 * time.Minute),
		Status:    entities.SessionStatusPending,
	}
	repo.sessions[validSession.Token] = validSession

	// Another expired session
	anotherExpiredSession := &entities.Session{
		Token:     "another-expired",
		CreatedAt: now.Add(-90 * time.Minute),
		ExpiresAt: now.Add(-60 * time.Minute),
		Status:    entities.SessionStatusActive,
	}
	repo.sessions[anotherExpiredSession.Token] = anotherExpiredSession

	// Run cleanup
	deletedCount, err := repo.CleanupExpiredSessions()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if deletedCount != 2 {
		t.Errorf("Expected to delete 2 sessions but deleted %d", deletedCount)
	}

	// Verify expired sessions are gone
	_, err = repo.GetSession("expired")
	if err == nil {
		t.Error("Expired session should be deleted")
	}

	_, err = repo.GetSession("another-expired")
	if err == nil {
		t.Error("Another expired session should be deleted")
	}

	// Verify valid session still exists
	_, err = repo.GetSession("valid")
	if err != nil {
		t.Errorf("Valid session should still exist: %v", err)
	}
}

func TestMemorySessionRepository_GetActiveSessionsCount(t *testing.T) {
	repo := NewMemorySessionRepository().(*MemorySessionRepository)

	// Initially should be 0
	count, err := repo.GetActiveSessionsCount()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 active sessions but got %d", count)
	}

	now := time.Now()

	// Add pending session (not active)
	pendingSession := &entities.Session{
		Token:     "pending",
		CreatedAt: now,
		ExpiresAt: now.Add(30 * time.Minute),
		Status:    entities.SessionStatusPending,
	}
	repo.sessions[pendingSession.Token] = pendingSession

	// Add active session
	activeSession := &entities.Session{
		Token:     "active",
		CreatedAt: now,
		ExpiresAt: now.Add(30 * time.Minute),
		Status:    entities.SessionStatusActive,
	}
	repo.sessions[activeSession.Token] = activeSession

	// Add expired active session (should not count)
	expiredActiveSession := &entities.Session{
		Token:     "expired-active",
		CreatedAt: now.Add(-60 * time.Minute),
		ExpiresAt: now.Add(-30 * time.Minute),
		Status:    entities.SessionStatusActive,
	}
	repo.sessions[expiredActiveSession.Token] = expiredActiveSession

	// Check count
	count, err = repo.GetActiveSessionsCount()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 active session but got %d", count)
	}
}
