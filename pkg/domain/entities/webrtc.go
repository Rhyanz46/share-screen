package entities

// WebRTCOffer represents a WebRTC offer
type WebRTCOffer struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

// WebRTCAnswer represents a WebRTC answer
type WebRTCAnswer struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

// IsValid checks if the WebRTC offer is valid
func (o *WebRTCOffer) IsValid() bool {
	return o != nil && o.Type != "" && o.SDP != ""
}

// IsValid checks if the WebRTC answer is valid
func (a *WebRTCAnswer) IsValid() bool {
	return a != nil && a.Type != "" && a.SDP != ""
}
