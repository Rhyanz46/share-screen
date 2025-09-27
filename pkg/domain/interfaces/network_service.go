package interfaces

// NetworkService defines the contract for network-related operations
type NetworkService interface {
	// GetLANIP returns the local area network IP address
	GetLANIP() string
}
