package dto

import "share-screen/pkg/domain/entities"

// CreateSessionResponse represents the response for creating a new session
type CreateSessionResponse struct {
	Token string `json:"token"`
}

// SubmitOfferRequest represents the request for submitting a WebRTC offer
type SubmitOfferRequest struct {
	Token string                `json:"token"`
	Offer *entities.WebRTCOffer `json:"sdp"`
}

// GetOfferRequest represents the request for getting a WebRTC offer
type GetOfferRequest struct {
	Token string `json:"token"`
}

// GetOfferResponse represents the response for getting a WebRTC offer
type GetOfferResponse struct {
	Offer *entities.WebRTCOffer `json:"offer"`
}

// SubmitAnswerRequest represents the request for submitting a WebRTC answer
type SubmitAnswerRequest struct {
	Token  string                 `json:"token"`
	Answer *entities.WebRTCAnswer `json:"sdp"`
}

// GetAnswerRequest represents the request for getting a WebRTC answer
type GetAnswerRequest struct {
	Token string `json:"token"`
}

// GetAnswerResponse represents the response for getting a WebRTC answer
type GetAnswerResponse struct {
	Answer *entities.WebRTCAnswer `json:"answer"`
}
