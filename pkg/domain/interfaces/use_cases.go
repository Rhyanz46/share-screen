package interfaces

import (
	"share-screen/pkg/domain/entities"
	"share-screen/pkg/usecase/dto"
)

// SessionUseCase defines the contract for session-related business logic
type SessionUseCase interface {
	// CreateSession creates a new screen sharing session
	CreateSession() (*dto.CreateSessionResponse, error)

	// SubmitOffer submits a WebRTC offer for a session
	SubmitOffer(request *dto.SubmitOfferRequest) error

	// GetOffer retrieves a WebRTC offer for a session
	GetOffer(request *dto.GetOfferRequest) (*dto.GetOfferResponse, error)

	// SubmitAnswer submits a WebRTC answer for a session
	SubmitAnswer(request *dto.SubmitAnswerRequest) error

	// GetAnswer retrieves a WebRTC answer for a session
	GetAnswer(request *dto.GetAnswerRequest) (*dto.GetAnswerResponse, error)
}

// ServerInfoUseCase defines the contract for server information
type ServerInfoUseCase interface {
	// GetServerInfo returns server information including network details
	GetServerInfo(host string) (*entities.ServerInfo, error)
}