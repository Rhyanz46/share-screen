package usecases

import (
	"testing"

	"share-screen/test/mocks"
)

func TestServerInfoUseCase_GetServerInfo(t *testing.T) {
	tests := []struct {
		name           string
		host           string
		mockLANIP      string
		stunServer     string
		version        string
		expectedHost   string
		expectedLANIP  string
		expectedSTUN   string
		expectedVersion string
	}{
		{
			name:           "successful server info retrieval",
			host:           "localhost:8080",
			mockLANIP:      "192.168.1.100",
			stunServer:     "stun:stun.l.google.com:19302",
			version:        "1.0.0",
			expectedHost:   "localhost:8080",
			expectedLANIP:  "192.168.1.100",
			expectedSTUN:   "stun:stun.l.google.com:19302",
			expectedVersion: "1.0.0",
		},
		{
			name:           "empty host",
			host:           "",
			mockLANIP:      "10.0.0.100",
			stunServer:     "stun:stun.example.com:3478",
			version:        "2.0.0",
			expectedHost:   "",
			expectedLANIP:  "10.0.0.100",
			expectedSTUN:   "stun:stun.example.com:3478",
			expectedVersion: "2.0.0",
		},
		{
			name:           "empty LAN IP",
			host:           "example.com:443",
			mockLANIP:      "",
			stunServer:     "stun:stun.l.google.com:19302",
			version:        "1.2.3",
			expectedHost:   "example.com:443",
			expectedLANIP:  "",
			expectedSTUN:   "stun:stun.l.google.com:19302",
			expectedVersion: "1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockNetworkService := mocks.NewMockNetworkService()
			mockNetworkService.SetLANIP(tt.mockLANIP)

			useCase := NewServerInfoUseCase(mockNetworkService, tt.stunServer, tt.version)

			// Execute
			result, err := useCase.GetServerInfo(tt.host)

			// Assert
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("Expected result but got nil")
			}

			if result.Host != tt.expectedHost {
				t.Errorf("Expected host %q but got %q", tt.expectedHost, result.Host)
			}

			if result.LANIP != tt.expectedLANIP {
				t.Errorf("Expected LAN IP %q but got %q", tt.expectedLANIP, result.LANIP)
			}

			if result.STUNServer != tt.expectedSTUN {
				t.Errorf("Expected STUN server %q but got %q", tt.expectedSTUN, result.STUNServer)
			}

			if result.Version != tt.expectedVersion {
				t.Errorf("Expected version %q but got %q", tt.expectedVersion, result.Version)
			}
		})
	}
}