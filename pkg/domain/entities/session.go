package entities

import (
	"time"
)

// Session represents a screen sharing session
type Session struct {
	Token     string
	Offer     *WebRTCOffer
	Answer    *WebRTCAnswer
	CreatedAt time.Time
	ExpiresAt time.Time
	Status    SessionStatus
}

// SessionStatus represents the current status of a session
type SessionStatus string

const (
	SessionStatusPending   SessionStatus = "pending"
	SessionStatusActive    SessionStatus = "active"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusExpired   SessionStatus = "expired"
)

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsActive checks if the session is currently active
func (s *Session) IsActive() bool {
	return s.Status == SessionStatusActive && !s.IsExpired()
}

// CanAcceptOffer checks if the session can accept a WebRTC offer
func (s *Session) CanAcceptOffer() bool {
	return s.Status == SessionStatusPending && !s.IsExpired()
}

// CanAcceptAnswer checks if the session can accept a WebRTC answer
func (s *Session) CanAcceptAnswer() bool {
	return s.Offer != nil && s.Answer == nil && !s.IsExpired()
}
