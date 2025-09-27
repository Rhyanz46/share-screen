package usecases

import (
	"testing"
	"time"

	"share-screen/pkg/domain/entities"
	"share-screen/pkg/usecase/dto"
	"share-screen/test/mocks"
)

func TestSessionUseCase_CreateSession(t *testing.T) {
	tests := []struct {
		name               string
		shouldFailCreate   bool
		expectedError      error
	}{
		{
			name:               "successful session creation",
			shouldFailCreate:   false,
			expectedError:      nil,
		},
		{
			name:               "failed session creation",
			shouldFailCreate:   true,
			expectedError:      nil, // We expect an error but don't check the specific type
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockSessionRepository()
			mockRepo.ShouldFailCreateSession = tt.shouldFailCreate

			useCase := NewSessionUseCase(mockRepo, 30*time.Minute)

			// Execute
			response, err := useCase.CreateSession()

			// Assert
			if tt.shouldFailCreate {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if response != nil {
					t.Error("Expected nil response on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if response == nil {
					t.Error("Expected response but got nil")
				}
				if response != nil && response.Token == "" {
					t.Error("Expected token but got empty string")
				}
			}
		})
	}
}

func TestSessionUseCase_SubmitOffer(t *testing.T) {
	tests := []struct {
		name            string
		request         *dto.SubmitOfferRequest
		setupSession    func(*mocks.MockSessionRepository)
		expectedError   error
	}{
		{
			name: "successful offer submission",
			request: &dto.SubmitOfferRequest{
				Token: "test-token",
				Offer: &entities.WebRTCOffer{
					Type: "offer",
					SDP:  "test-sdp",
				},
			},
			setupSession: func(repo *mocks.MockSessionRepository) {
				session := &entities.Session{
					Token:     "test-token",
					CreatedAt: time.Now(),
					ExpiresAt: time.Now().Add(30 * time.Minute),
					Status:    entities.SessionStatusPending,
				}
				repo.SetSession(session)
			},
			expectedError: nil,
		},
		{
			name: "invalid offer - nil",
			request: &dto.SubmitOfferRequest{
				Token: "test-token",
				Offer: nil,
			},
			setupSession: func(repo *mocks.MockSessionRepository) {},
			expectedError: ErrInvalidOffer,
		},
		{
			name: "invalid offer - empty type",
			request: &dto.SubmitOfferRequest{
				Token: "test-token",
				Offer: &entities.WebRTCOffer{
					Type: "",
					SDP:  "test-sdp",
				},
			},
			setupSession: func(repo *mocks.MockSessionRepository) {},
			expectedError: ErrInvalidOffer,
		},
		{
			name: "session not found",
			request: &dto.SubmitOfferRequest{
				Token: "non-existent-token",
				Offer: &entities.WebRTCOffer{
					Type: "offer",
					SDP:  "test-sdp",
				},
			},
			setupSession: func(repo *mocks.MockSessionRepository) {},
			expectedError: ErrSessionNotFound,
		},
		{
			name: "expired session",
			request: &dto.SubmitOfferRequest{
				Token: "expired-token",
				Offer: &entities.WebRTCOffer{
					Type: "offer",
					SDP:  "test-sdp",
				},
			},
			setupSession: func(repo *mocks.MockSessionRepository) {
				session := &entities.Session{
					Token:     "expired-token",
					CreatedAt: time.Now().Add(-60 * time.Minute),
					ExpiresAt: time.Now().Add(-30 * time.Minute),
					Status:    entities.SessionStatusPending,
				}
				repo.SetSession(session)
			},
			expectedError: ErrSessionExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockSessionRepository()
			tt.setupSession(mockRepo)

			useCase := NewSessionUseCase(mockRepo, 30*time.Minute)

			// Execute
			err := useCase.SubmitOffer(tt.request)

			// Assert
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v but got none", tt.expectedError)
				}
				if err != tt.expectedError {
					t.Errorf("Expected error %v but got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSessionUseCase_GetOffer(t *testing.T) {
	tests := []struct {
		name            string
		request         *dto.GetOfferRequest
		setupSession    func(*mocks.MockSessionRepository)
		expectedError   error
		shouldHaveOffer bool
	}{
		{
			name: "successful offer retrieval",
			request: &dto.GetOfferRequest{
				Token: "test-token",
			},
			setupSession: func(repo *mocks.MockSessionRepository) {
				session := &entities.Session{
					Token:     "test-token",
					CreatedAt: time.Now(),
					ExpiresAt: time.Now().Add(30 * time.Minute),
					Status:    entities.SessionStatusActive,
					Offer: &entities.WebRTCOffer{
						Type: "offer",
						SDP:  "test-sdp",
					},
				}
				repo.SetSession(session)
			},
			expectedError:   nil,
			shouldHaveOffer: true,
		},
		{
			name: "session not found",
			request: &dto.GetOfferRequest{
				Token: "non-existent-token",
			},
			setupSession:    func(repo *mocks.MockSessionRepository) {},
			expectedError:   ErrSessionNotFound,
			shouldHaveOffer: false,
		},
		{
			name: "offer not found",
			request: &dto.GetOfferRequest{
				Token: "test-token",
			},
			setupSession: func(repo *mocks.MockSessionRepository) {
				session := &entities.Session{
					Token:     "test-token",
					CreatedAt: time.Now(),
					ExpiresAt: time.Now().Add(30 * time.Minute),
					Status:    entities.SessionStatusPending,
					Offer:     nil,
				}
				repo.SetSession(session)
			},
			expectedError:   ErrOfferNotFound,
			shouldHaveOffer: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockSessionRepository()
			tt.setupSession(mockRepo)

			useCase := NewSessionUseCase(mockRepo, 30*time.Minute)

			// Execute
			response, err := useCase.GetOffer(tt.request)

			// Assert
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v but got none", tt.expectedError)
				}
				if err != tt.expectedError {
					t.Errorf("Expected error %v but got %v", tt.expectedError, err)
				}
				if response != nil {
					t.Error("Expected nil response on error")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.shouldHaveOffer {
					if response == nil {
						t.Error("Expected response but got nil")
					}
					if response != nil && response.Offer == nil {
						t.Error("Expected offer but got nil")
					}
				}
			}
		})
	}
}

func TestSessionUseCase_SubmitAnswer(t *testing.T) {
	tests := []struct {
		name            string
		request         *dto.SubmitAnswerRequest
		setupSession    func(*mocks.MockSessionRepository)
		expectedError   error
	}{
		{
			name: "successful answer submission",
			request: &dto.SubmitAnswerRequest{
				Token: "test-token",
				Answer: &entities.WebRTCAnswer{
					Type: "answer",
					SDP:  "test-answer-sdp",
				},
			},
			setupSession: func(repo *mocks.MockSessionRepository) {
				session := &entities.Session{
					Token:     "test-token",
					CreatedAt: time.Now(),
					ExpiresAt: time.Now().Add(30 * time.Minute),
					Status:    entities.SessionStatusActive,
					Offer: &entities.WebRTCOffer{
						Type: "offer",
						SDP:  "test-sdp",
					},
					Answer: nil,
				}
				repo.SetSession(session)
			},
			expectedError: nil,
		},
		{
			name: "invalid answer - nil",
			request: &dto.SubmitAnswerRequest{
				Token:  "test-token",
				Answer: nil,
			},
			setupSession:  func(repo *mocks.MockSessionRepository) {},
			expectedError: ErrInvalidAnswer,
		},
		{
			name: "answer already exists",
			request: &dto.SubmitAnswerRequest{
				Token: "test-token",
				Answer: &entities.WebRTCAnswer{
					Type: "answer",
					SDP:  "test-answer-sdp",
				},
			},
			setupSession: func(repo *mocks.MockSessionRepository) {
				session := &entities.Session{
					Token:     "test-token",
					CreatedAt: time.Now(),
					ExpiresAt: time.Now().Add(30 * time.Minute),
					Status:    entities.SessionStatusActive,
					Offer: &entities.WebRTCOffer{
						Type: "offer",
						SDP:  "test-sdp",
					},
					Answer: &entities.WebRTCAnswer{
						Type: "answer",
						SDP:  "existing-answer-sdp",
					},
				}
				repo.SetSession(session)
			},
			expectedError: ErrAnswerAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockSessionRepository()
			tt.setupSession(mockRepo)

			useCase := NewSessionUseCase(mockRepo, 30*time.Minute)

			// Execute
			err := useCase.SubmitAnswer(tt.request)

			// Assert
			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v but got none", tt.expectedError)
				}
				if err != tt.expectedError {
					t.Errorf("Expected error %v but got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}