package network

import (
	"net"
	"testing"
)

func TestNetworkService_GetLANIP(t *testing.T) {
	service := NewNetworkService().(*NetworkService)

	ip := service.GetLANIP()

	// The actual IP will depend on the test environment
	// We just test that the function doesn't panic and returns a string
	if ip == "" {
		// This might happen in some test environments
		t.Log("No LAN IP found (this might be normal in test environment)")
	} else {
		// If an IP is returned, it should be a valid IPv4 address
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			t.Errorf("Returned IP %q is not a valid IP address", ip)
		}

		// Convert to IPv4 to ensure it's IPv4
		ipv4 := parsedIP.To4()
		if ipv4 == nil {
			t.Errorf("Returned IP %q is not a valid IPv4 address", ip)
		}

		// Check if it's in private ranges (this is what the function should return)
		isPrivate := false
		if ipv4[0] == 10 {
			isPrivate = true
		} else if ipv4[0] == 192 && ipv4[1] == 168 {
			isPrivate = true
		} else if ipv4[0] == 172 && ipv4[1] >= 16 && ipv4[1] <= 31 {
			isPrivate = true
		}

		if !isPrivate {
			t.Logf("Warning: Returned IP %q is not in typical private ranges", ip)
		}
	}
}

func TestNetworkService_GetLANIP_Integration(t *testing.T) {
	// This is more of an integration test that verifies the function
	// works with the actual network interfaces
	service := NewNetworkService()

	// Call multiple times to ensure consistency
	ip1 := service.GetLANIP()
	ip2 := service.GetLANIP()

	if ip1 != ip2 {
		t.Errorf("GetLANIP should return consistent results: got %q and %q", ip1, ip2)
	}
}
