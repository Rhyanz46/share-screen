package entities

// ServerInfo represents server information
type ServerInfo struct {
	Host       string `json:"host"`
	LANIP      string `json:"lanIP"`
	STUNServer string `json:"stunServer,omitempty"`
	Version    string `json:"version,omitempty"`
}
