package entities

import (
	"testing"
	"time"
)

func TestSession_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		session  *Session
		expected bool
	}{
		{
			name: "session not expired",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-10 * time.Minute),
				ExpiresAt: time.Now().Add(10 * time.Minute),
				Status:    SessionStatusPending,
			},
			expected: false,
		},
		{
			name: "session expired",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-30 * time.Minute),
				ExpiresAt: time.Now().Add(-10 * time.Minute),
				Status:    SessionStatusPending,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.IsExpired()
			if result != tt.expected {
				t.Errorf("IsExpired() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSession_IsActive(t *testing.T) {
	tests := []struct {
		name     string
		session  *Session
		expected bool
	}{
		{
			name: "active session",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-10 * time.Minute),
				ExpiresAt: time.Now().Add(10 * time.Minute),
				Status:    SessionStatusActive,
			},
			expected: true,
		},
		{
			name: "pending session",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-10 * time.Minute),
				ExpiresAt: time.Now().Add(10 * time.Minute),
				Status:    SessionStatusPending,
			},
			expected: false,
		},
		{
			name: "expired active session",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-30 * time.Minute),
				ExpiresAt: time.Now().Add(-10 * time.Minute),
				Status:    SessionStatusActive,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.IsActive()
			if result != tt.expected {
				t.Errorf("IsActive() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSession_CanAcceptOffer(t *testing.T) {
	tests := []struct {
		name     string
		session  *Session
		expected bool
	}{
		{
			name: "pending session can accept offer",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-10 * time.Minute),
				ExpiresAt: time.Now().Add(10 * time.Minute),
				Status:    SessionStatusPending,
			},
			expected: true,
		},
		{
			name: "active session cannot accept offer",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-10 * time.Minute),
				ExpiresAt: time.Now().Add(10 * time.Minute),
				Status:    SessionStatusActive,
			},
			expected: false,
		},
		{
			name: "expired pending session cannot accept offer",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-30 * time.Minute),
				ExpiresAt: time.Now().Add(-10 * time.Minute),
				Status:    SessionStatusPending,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.CanAcceptOffer()
			if result != tt.expected {
				t.Errorf("CanAcceptOffer() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSession_CanAcceptAnswer(t *testing.T) {
	tests := []struct {
		name     string
		session  *Session
		expected bool
	}{
		{
			name: "session with offer can accept answer",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-10 * time.Minute),
				ExpiresAt: time.Now().Add(10 * time.Minute),
				Status:    SessionStatusActive,
				Offer:     &WebRTCOffer{Type: "offer", SDP: "test-sdp"},
				Answer:    nil,
			},
			expected: true,
		},
		{
			name: "session without offer cannot accept answer",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-10 * time.Minute),
				ExpiresAt: time.Now().Add(10 * time.Minute),
				Status:    SessionStatusActive,
				Offer:     nil,
				Answer:    nil,
			},
			expected: false,
		},
		{
			name: "session with existing answer cannot accept answer",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-10 * time.Minute),
				ExpiresAt: time.Now().Add(10 * time.Minute),
				Status:    SessionStatusActive,
				Offer:     &WebRTCOffer{Type: "offer", SDP: "test-sdp"},
				Answer:    &WebRTCAnswer{Type: "answer", SDP: "test-answer-sdp"},
			},
			expected: false,
		},
		{
			name: "expired session cannot accept answer",
			session: &Session{
				Token:     "test-token",
				CreatedAt: time.Now().Add(-30 * time.Minute),
				ExpiresAt: time.Now().Add(-10 * time.Minute),
				Status:    SessionStatusActive,
				Offer:     &WebRTCOffer{Type: "offer", SDP: "test-sdp"},
				Answer:    nil,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.session.CanAcceptAnswer()
			if result != tt.expected {
				t.Errorf("CanAcceptAnswer() = %v, want %v", result, tt.expected)
			}
		})
	}
}
