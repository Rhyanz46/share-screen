package network

import (
	"log"
	"net"

	"share-screen/pkg/domain/interfaces"
)

// NetworkService implements the NetworkService interface
type NetworkService struct{}

// NewNetworkService creates a new network service
func NewNetworkService() interfaces.NetworkService {
	return &NetworkService{}
}

// GetLANIP returns the local area network IP address
func (s *NetworkService) GetLANIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Error getting network interfaces: %v", err)
		return ""
	}
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP == nil || ipnet.IP.IsLoopback() {
				continue
			}
			ipv4 := ipnet.IP.To4()
			if ipv4 == nil {
				continue
			}
			// pick typical private ranges
			if ipv4[0] == 10 || (ipv4[0] == 192 && ipv4[1] == 168) || (ipv4[0] == 172 && ipv4[1] >= 16 && ipv4[1] <= 31) {
				return ipv4.String()
			}
		}
	}
	return ""
}