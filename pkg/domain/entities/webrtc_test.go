package entities

import "testing"

func TestWebRTCOffer_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		offer    *WebRTCOffer
		expected bool
	}{
		{
			name: "valid offer",
			offer: &WebRTCOffer{
				Type: "offer",
				SDP:  "v=0\no=- 123456789 123456789 IN IP4 192.168.1.1\n",
			},
			expected: true,
		},
		{
			name:     "nil offer",
			offer:    nil,
			expected: false,
		},
		{
			name: "empty type",
			offer: &WebRTCOffer{
				Type: "",
				SDP:  "v=0\no=- 123456789 123456789 IN IP4 192.168.1.1\n",
			},
			expected: false,
		},
		{
			name: "empty SDP",
			offer: &WebRTCOffer{
				Type: "offer",
				SDP:  "",
			},
			expected: false,
		},
		{
			name: "both empty",
			offer: &WebRTCOffer{
				Type: "",
				SDP:  "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.offer.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWebRTCAnswer_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		answer   *WebRTCAnswer
		expected bool
	}{
		{
			name: "valid answer",
			answer: &WebRTCAnswer{
				Type: "answer",
				SDP:  "v=0\no=- 123456789 123456789 IN IP4 192.168.1.1\n",
			},
			expected: true,
		},
		{
			name:     "nil answer",
			answer:   nil,
			expected: false,
		},
		{
			name: "empty type",
			answer: &WebRTCAnswer{
				Type: "",
				SDP:  "v=0\no=- 123456789 123456789 IN IP4 192.168.1.1\n",
			},
			expected: false,
		},
		{
			name: "empty SDP",
			answer: &WebRTCAnswer{
				Type: "answer",
				SDP:  "",
			},
			expected: false,
		},
		{
			name: "both empty",
			answer: &WebRTCAnswer{
				Type: "",
				SDP:  "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.answer.IsValid()
			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}