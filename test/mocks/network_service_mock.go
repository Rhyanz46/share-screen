package mocks

// MockNetworkService is a mock implementation of NetworkService interface
type MockNetworkService struct {
	LANIPToReturn string
}

// NewMockNetworkService creates a new mock network service
func NewMockNetworkService() *MockNetworkService {
	return &MockNetworkService{
		LANIPToReturn: "192.168.1.100", // Default mock IP
	}
}

// GetLANIP returns the configured mock LAN IP
func (m *MockNetworkService) GetLANIP() string {
	return m.LANIPToReturn
}

// SetLANIP sets the LAN IP to be returned (for testing purposes)
func (m *MockNetworkService) SetLANIP(ip string) {
	m.LANIPToReturn = ip
}