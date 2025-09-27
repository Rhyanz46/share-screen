package repository

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"sync"
	"time"

	"share-screen/pkg/domain/entities"
	"share-screen/pkg/domain/interfaces"
)

// MemorySessionRepository implements SessionRepository using in-memory storage
type MemorySessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]*entities.Session
}

// NewMemorySessionRepository creates a new in-memory session repository
func NewMemorySessionRepository() interfaces.SessionRepository {
	return &MemorySessionRepository{
		sessions: make(map[string]*entities.Session),
	}
}

// CreateSession creates a new session with a unique token
func (r *MemorySessionRepository) CreateSession(expiryDuration time.Duration) (*entities.Session, error) {
	token, err := r.generateToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	session := &entities.Session{
		Token:     token,
		CreatedAt: now,
		ExpiresAt: now.Add(expiryDuration),
		Status:    entities.SessionStatusPending,
	}

	r.mu.Lock()
	r.sessions[token] = session
	r.mu.Unlock()

	return session, nil
}

// GetSession retrieves a session by token
func (r *MemorySessionRepository) GetSession(token string) (*entities.Session, error) {
	r.mu.RLock()
	session, exists := r.sessions[token]
	r.mu.RUnlock()

	if !exists {
		return nil, ErrSessionNotFound
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
func (r *MemorySessionRepository) UpdateSession(session *entities.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.sessions[session.Token]
	if !exists {
		return ErrSessionNotFound
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

	r.sessions[session.Token] = &sessionCopy
	return nil
}

// DeleteSession removes a session
func (r *MemorySessionRepository) DeleteSession(token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, token)
	return nil
}

// CleanupExpiredSessions removes all expired sessions
func (r *MemorySessionRepository) CleanupExpiredSessions() (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var expiredTokens []string
	for token, session := range r.sessions {
		if session.IsExpired() {
			expiredTokens = append(expiredTokens, token)
		}
	}

	for _, token := range expiredTokens {
		delete(r.sessions, token)
	}

	if len(expiredTokens) > 0 {
		// Convert to truncated tokens for logging
		var truncatedTokens []string
		for _, token := range expiredTokens {
			if len(token) > 8 {
				truncatedTokens = append(truncatedTokens, token[:8]+"...")
			} else {
				truncatedTokens = append(truncatedTokens, token+"...")
			}
		}
		activeCount := len(r.sessions)
		log.Printf("🗑️  GC: cleaned up %d expired tokens: %v (active: %d)",
			len(expiredTokens), truncatedTokens, activeCount)
	}

	return len(expiredTokens), nil
}

// GetActiveSessionsCount returns the number of active sessions
func (r *MemorySessionRepository) GetActiveSessionsCount() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, session := range r.sessions {
		if session.IsActive() {
			count++
		}
	}

	return count, nil
}

// generateToken generates a random token for sessions
func (r *MemorySessionRepository) generateToken() (string, error) {
	b := make([]byte, 9)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(b)
	if len(token) > 8 {
		log.Printf("🆕 New token generated: %s...", token[:8])
	} else {
		log.Printf("🆕 New token generated: %s...", token)
	}
	return token, nil
}

// Define the error interface for the repository layer
var ErrSessionNotFound = &RepositoryError{Message: "session not found"}

type RepositoryError struct {
	Message string
}

func (e *RepositoryError) Error() string {
	return e.Message
}
