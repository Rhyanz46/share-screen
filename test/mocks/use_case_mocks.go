package mocks

import (
	"errors"

	"share-screen/pkg/domain/entities"
	"share-screen/pkg/usecase/dto"
)

// MockSessionUseCase is a mock implementation of SessionUseCase interface
type MockSessionUseCase struct {
	// For controlling behavior in tests
	ShouldFailCreateSession bool
	ShouldFailSubmitOffer   bool
	ShouldFailGetOffer      bool
	ShouldFailSubmitAnswer  bool
	ShouldFailGetAnswer     bool

	// For returning specific data
	CreateSessionResponse *dto.CreateSessionResponse
	GetOfferResponse      *dto.GetOfferResponse
	GetAnswerResponse     *dto.GetAnswerResponse
}

// NewMockSessionUseCase creates a new mock session use case
func NewMockSessionUseCase() *MockSessionUseCase {
	return &MockSessionUseCase{
		CreateSessionResponse: &dto.CreateSessionResponse{Token: "mock-token"},
		GetOfferResponse: &dto.GetOfferResponse{
			Offer: &entities.WebRTCOffer{Type: "offer", SDP: "mock-sdp"},
		},
		GetAnswerResponse: &dto.GetAnswerResponse{
			Answer: &entities.WebRTCAnswer{Type: "answer", SDP: "mock-answer-sdp"},
		},
	}
}

// CreateSession creates a new screen sharing session
func (m *MockSessionUseCase) CreateSession() (*dto.CreateSessionResponse, error) {
	if m.ShouldFailCreateSession {
		return nil, errors.New("mock create session error")
	}
	return m.CreateSessionResponse, nil
}

// SubmitOffer submits a WebRTC offer for a session
func (m *MockSessionUseCase) SubmitOffer(request *dto.SubmitOfferRequest) error {
	if m.ShouldFailSubmitOffer {
		return errors.New("mock submit offer error")
	}
	return nil
}

// GetOffer retrieves a WebRTC offer for a session
func (m *MockSessionUseCase) GetOffer(request *dto.GetOfferRequest) (*dto.GetOfferResponse, error) {
	if m.ShouldFailGetOffer {
		return nil, errors.New("mock get offer error")
	}
	return m.GetOfferResponse, nil
}

// SubmitAnswer submits a WebRTC answer for a session
func (m *MockSessionUseCase) SubmitAnswer(request *dto.SubmitAnswerRequest) error {
	if m.ShouldFailSubmitAnswer {
		return errors.New("mock submit answer error")
	}
	return nil
}

// GetAnswer retrieves a WebRTC answer for a session
func (m *MockSessionUseCase) GetAnswer(request *dto.GetAnswerRequest) (*dto.GetAnswerResponse, error) {
	if m.ShouldFailGetAnswer {
		return nil, errors.New("mock get answer error")
	}
	return m.GetAnswerResponse, nil
}

// MockServerInfoUseCase is a mock implementation of ServerInfoUseCase interface
type MockServerInfoUseCase struct {
	// For controlling behavior in tests
	ShouldFailGetServerInfo bool

	// For returning specific data
	ServerInfo *entities.ServerInfo
}

// NewMockServerInfoUseCase creates a new mock server info use case
func NewMockServerInfoUseCase() *MockServerInfoUseCase {
	return &MockServerInfoUseCase{
		ServerInfo: &entities.ServerInfo{
			Host:       "mock-host",
			LANIP:      "192.168.1.100",
			STUNServer: "stun:mock.com:19302",
			Version:    "1.0.0",
		},
	}
}

// GetServerInfo returns server information including network details
func (m *MockServerInfoUseCase) GetServerInfo(host string) (*entities.ServerInfo, error) {
	if m.ShouldFailGetServerInfo {
		return nil, errors.New("mock get server info error")
	}

	// Update host with the provided value
	result := *m.ServerInfo
	result.Host = host
	return &result, nil
}