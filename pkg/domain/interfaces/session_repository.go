package interfaces

import (
	"time"

	"share-screen/pkg/domain/entities"
)

// SessionRepository defines the contract for session data storage
type SessionRepository interface {
	// CreateSession creates a new session with a unique token
	CreateSession(expiryDuration time.Duration) (*entities.Session, error)

	// GetSession retrieves a session by token
	GetSession(token string) (*entities.Session, error)

	// UpdateSession updates an existing session
	UpdateSession(session *entities.Session) error

	// DeleteSession removes a session
	DeleteSession(token string) error

	// CleanupExpiredSessions removes all expired sessions
	CleanupExpiredSessions() (int, error)

	// GetActiveSessionsCount returns the number of active sessions
	GetActiveSessionsCount() (int, error)
}
