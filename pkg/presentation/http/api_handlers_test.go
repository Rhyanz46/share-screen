package http

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"share-screen/pkg/domain/entities"
	"share-screen/pkg/usecase/dto"
	"share-screen/test/mocks"
)

func TestAPIHandlers_HandleNewToken(t *testing.T) {
	tests := []struct {
		name                 string
		method               string
		shouldFailCreate     bool
		expectedStatusCode   int
		expectTokenInResponse bool
	}{
		{
			name:                 "successful token creation",
			method:               "POST",
			shouldFailCreate:     false,
			expectedStatusCode:   200,
			expectTokenInResponse: true,
		},
		{
			name:                 "method not allowed",
			method:               "GET",
			shouldFailCreate:     false,
			expectedStatusCode:   405,
			expectTokenInResponse: false,
		},
		{
			name:                 "failed token creation",
			method:               "POST",
			shouldFailCreate:     true,
			expectedStatusCode:   500,
			expectTokenInResponse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSessionUseCase := mocks.NewMockSessionUseCase()
			mockServerInfoUseCase := mocks.NewMockServerInfoUseCase()
			mockSessionUseCase.ShouldFailCreateSession = tt.shouldFailCreate

			handlers := NewAPIHandlers(mockSessionUseCase, mockServerInfoUseCase)

			// Create request
			req := httptest.NewRequest(tt.method, "/api/new", nil)
			w := httptest.NewRecorder()

			// Execute
			handlers.HandleNewToken(w, req)

			// Assert
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectTokenInResponse {
				var response dto.CreateSessionResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.Token == "" {
					t.Error("Expected token in response but got empty string")
				}
			}
		})
	}
}

func TestAPIHandlers_HandleOffer_POST(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		shouldFailSubmit   bool
		expectedStatusCode int
	}{
		{
			name: "successful offer submission",
			requestBody: dto.SubmitOfferRequest{
				Token: "test-token",
				Offer: &entities.WebRTCOffer{
					Type: "offer",
					SDP:  "test-sdp",
				},
			},
			shouldFailSubmit:   false,
			expectedStatusCode: 204,
		},
		{
			name:               "invalid JSON",
			requestBody:        "invalid-json",
			shouldFailSubmit:   false,
			expectedStatusCode: 400,
		},
		{
			name: "failed offer submission",
			requestBody: dto.SubmitOfferRequest{
				Token: "test-token",
				Offer: &entities.WebRTCOffer{
					Type: "offer",
					SDP:  "test-sdp",
				},
			},
			shouldFailSubmit:   true,
			expectedStatusCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSessionUseCase := mocks.NewMockSessionUseCase()
			mockServerInfoUseCase := mocks.NewMockServerInfoUseCase()
			mockSessionUseCase.ShouldFailSubmitOffer = tt.shouldFailSubmit

			handlers := NewAPIHandlers(mockSessionUseCase, mockServerInfoUseCase)

			// Create request body
			var bodyBytes []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest("POST", "/api/offer", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			handlers.HandleOffer(w, req)

			// Assert
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", tt.expectedStatusCode, w.Code)
			}
		})
	}
}

func TestAPIHandlers_HandleOffer_GET(t *testing.T) {
	tests := []struct {
		name               string
		token              string
		shouldFailGet      bool
		expectedStatusCode int
		expectOfferInResponse bool
	}{
		{
			name:               "successful offer retrieval",
			token:              "test-token",
			shouldFailGet:      false,
			expectedStatusCode: 200,
			expectOfferInResponse: true,
		},
		{
			name:               "failed offer retrieval",
			token:              "test-token",
			shouldFailGet:      true,
			expectedStatusCode: 500,
			expectOfferInResponse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSessionUseCase := mocks.NewMockSessionUseCase()
			mockServerInfoUseCase := mocks.NewMockServerInfoUseCase()
			mockSessionUseCase.ShouldFailGetOffer = tt.shouldFailGet

			handlers := NewAPIHandlers(mockSessionUseCase, mockServerInfoUseCase)

			// Create request
			req := httptest.NewRequest("GET", "/api/offer?token="+tt.token, nil)
			w := httptest.NewRecorder()

			// Execute
			handlers.HandleOffer(w, req)

			// Assert
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectOfferInResponse {
				var offer entities.WebRTCOffer
				err := json.Unmarshal(w.Body.Bytes(), &offer)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if offer.Type == "" {
					t.Error("Expected offer type in response but got empty string")
				}
			}
		})
	}
}

func TestAPIHandlers_HandleAnswer_POST(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		shouldFailSubmit   bool
		expectedStatusCode int
	}{
		{
			name: "successful answer submission",
			requestBody: dto.SubmitAnswerRequest{
				Token: "test-token",
				Answer: &entities.WebRTCAnswer{
					Type: "answer",
					SDP:  "test-answer-sdp",
				},
			},
			shouldFailSubmit:   false,
			expectedStatusCode: 204,
		},
		{
			name:               "invalid JSON",
			requestBody:        "invalid-json",
			shouldFailSubmit:   false,
			expectedStatusCode: 400,
		},
		{
			name: "failed answer submission",
			requestBody: dto.SubmitAnswerRequest{
				Token: "test-token",
				Answer: &entities.WebRTCAnswer{
					Type: "answer",
					SDP:  "test-answer-sdp",
				},
			},
			shouldFailSubmit:   true,
			expectedStatusCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSessionUseCase := mocks.NewMockSessionUseCase()
			mockServerInfoUseCase := mocks.NewMockServerInfoUseCase()
			mockSessionUseCase.ShouldFailSubmitAnswer = tt.shouldFailSubmit

			handlers := NewAPIHandlers(mockSessionUseCase, mockServerInfoUseCase)

			// Create request body
			var bodyBytes []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest("POST", "/api/answer", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			handlers.HandleAnswer(w, req)

			// Assert
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", tt.expectedStatusCode, w.Code)
			}
		})
	}
}

func TestAPIHandlers_HandleInfo(t *testing.T) {
	tests := []struct {
		name               string
		host               string
		shouldFailGet      bool
		expectedStatusCode int
		expectInfoInResponse bool
	}{
		{
			name:               "successful info retrieval",
			host:               "localhost:8080",
			shouldFailGet:      false,
			expectedStatusCode: 200,
			expectInfoInResponse: true,
		},
		{
			name:               "failed info retrieval",
			host:               "localhost:8080",
			shouldFailGet:      true,
			expectedStatusCode: 500,
			expectInfoInResponse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSessionUseCase := mocks.NewMockSessionUseCase()
			mockServerInfoUseCase := mocks.NewMockServerInfoUseCase()
			mockServerInfoUseCase.ShouldFailGetServerInfo = tt.shouldFailGet

			handlers := NewAPIHandlers(mockSessionUseCase, mockServerInfoUseCase)

			// Create request
			req := httptest.NewRequest("GET", "/api/info", nil)
			req.Host = tt.host
			w := httptest.NewRecorder()

			// Execute
			handlers.HandleInfo(w, req)

			// Assert
			if w.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", tt.expectedStatusCode, w.Code)
			}

			if tt.expectInfoInResponse {
				var info entities.ServerInfo
				err := json.Unmarshal(w.Body.Bytes(), &info)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if info.Host != tt.host {
					t.Errorf("Expected host %q but got %q", tt.host, info.Host)
				}
			}
		})
	}
}

func TestAPIHandlers_HandleOffer_MethodNotAllowed(t *testing.T) {
	mockSessionUseCase := mocks.NewMockSessionUseCase()
	mockServerInfoUseCase := mocks.NewMockServerInfoUseCase()
	handlers := NewAPIHandlers(mockSessionUseCase, mockServerInfoUseCase)

	req := httptest.NewRequest("DELETE", "/api/offer", nil)
	w := httptest.NewRecorder()

	handlers.HandleOffer(w, req)

	if w.Code != 405 {
		t.Errorf("Expected status code 405 but got %d", w.Code)
	}
}