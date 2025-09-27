package mocks

import (
	"time"

	"share-screen/pkg/domain/entities"
)

// MockSessionRepository is a mock implementation of SessionRepository interface
type MockSessionRepository struct {
	sessions map[string]*entities.Session

	// For controlling behavior in tests
	ShouldFailCreateSession bool
	ShouldFailUpdateSession bool
	ShouldFailGetSession    bool
}

// NewMockSessionRepository creates a new mock session repository
func NewMockSessionRepository() *MockSessionRepository {
	return &MockSessionRepository{
		sessions: make(map[string]*entities.Session),
	}
}

// CreateSession creates a new session with a unique token
func (m *MockSessionRepository) CreateSession(expiryDuration time.Duration) (*entities.Session, error) {
	if m.ShouldFailCreateSession {
		return nil, mockError("failed to create session")
	}

	token := "mock-token-" + time.Now().Format("150405")
	now := time.Now()
	session := &entities.Session{
		Token:     token,
		CreatedAt: now,
		ExpiresAt: now.Add(expiryDuration),
		Status:    entities.SessionStatusPending,
	}

	m.sessions[token] = session
	return session, nil
}

// GetSession retrieves a session by token
func (m *MockSessionRepository) GetSession(token string) (*entities.Session, error) {
	if m.ShouldFailGetSession {
		return nil, mockError("failed to get session")
	}

	session, exists := m.sessions[token]
	if !exists {
		return nil, mockError("session not found")
	}

	// Return a copy to prevent external modifications
	sessionCopy := *session
	if session.Offer != nil {
		offerCopy := *session.Offer
		sessionCopy.Offer = &offerCopy
	}
	if session.Answer != nil {
		answerCopy := *session.Answer
		sessionCopy.Answer = &answerCopy
	}

	return &sessionCopy, nil
}

// UpdateSession updates an existing session
func (m *MockSessionRepository) UpdateSession(session *entities.Session) error {
	if m.ShouldFailUpdateSession {
		return mockError("failed to update session")
	}

	_, exists := m.sessions[session.Token]
	if !exists {
		return mockError("session not found")
	}

	// Create a copy to store
	sessionCopy := *session
	if session.Offer != nil {
		offerCopy := *session.Offer
		sessionCopy.Offer = &offerCopy
	}
	if session.Answer != nil {
		answerCopy := *session.Answer
		sessionCopy.Answer = &answerCopy
	}

	m.sessions[session.Token] = &sessionCopy
	return nil
}

// DeleteSession removes a session
func (m *MockSessionRepository) DeleteSession(token string) error {
	delete(m.sessions, token)
	return nil
}

// CleanupExpiredSessions removes all expired sessions
func (m *MockSessionRepository) CleanupExpiredSessions() (int, error) {
	var expiredTokens []string
	for token, session := range m.sessions {
		if session.IsExpired() {
			expiredTokens = append(expiredTokens, token)
		}
	}

	for _, token := range expiredTokens {
		delete(m.sessions, token)
	}

	return len(expiredTokens), nil
}

// GetActiveSessionsCount returns the number of active sessions
func (m *MockSessionRepository) GetActiveSessionsCount() (int, error) {
	count := 0
	for _, session := range m.sessions {
		if session.IsActive() {
			count++
		}
	}
	return count, nil
}

// SetSession directly sets a session (for testing purposes)
func (m *MockSessionRepository) SetSession(session *entities.Session) {
	m.sessions[session.Token] = session
}

// GetSessionCount returns the total number of sessions (for testing purposes)
func (m *MockSessionRepository) GetSessionCount() int {
	return len(m.sessions)
}

// Clear removes all sessions (for testing purposes)
func (m *MockSessionRepository) Clear() {
	m.sessions = make(map[string]*entities.Session)
}

// mockError creates a simple error for testing
type mockError string

func (e mockError) Error() string {
	return string(e)
}