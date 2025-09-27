package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"share-screen/pkg/domain/entities"
	"share-screen/pkg/infrastructure/network"
	"share-screen/pkg/infrastructure/repository"
	httphandlers "share-screen/pkg/presentation/http"
	"share-screen/pkg/usecase/dto"
	"share-screen/pkg/usecase/usecases"
)

// TestHTTPAPIIntegration tests the complete HTTP API flow
func TestHTTPAPIIntegration(t *testing.T) {
	// Setup real dependencies
	sessionRepo := repository.NewMemorySessionRepository()
	networkService := network.NewNetworkService()

	sessionUseCase := usecases.NewSessionUseCase(sessionRepo, 30*time.Minute)
	serverInfoUseCase := usecases.NewServerInfoUseCase(networkService, "stun:test.com:19302", "1.0.0")

	apiHandlers := httphandlers.NewAPIHandlers(sessionUseCase, serverInfoUseCase)

	t.Run("complete HTTP API workflow", func(t *testing.T) {
		// Step 1: Create a new session
		req := httptest.NewRequest("POST", "/api/new", nil)
		w := httptest.NewRecorder()

		apiHandlers.HandleNewToken(w, req)

		if w.Code != 200 {
			t.Fatalf("Expected status 200 but got %d", w.Code)
		}

		var createResponse dto.CreateSessionResponse
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		if err != nil {
			t.Fatalf("Failed to unmarshal create response: %v", err)
		}

		if createResponse.Token == "" {
			t.Fatal("Expected token but got empty string")
		}

		token := createResponse.Token

		// Step 2: Submit an offer via HTTP
		offerRequest := dto.SubmitOfferRequest{
			Token: token,
			Offer: &entities.WebRTCOffer{
				Type: "offer",
				SDP:  "v=0\no=- 123456789 123456789 IN IP4 192.168.1.1\ns=-\nt=0 0\n",
			},
		}

		offerBody, err := json.Marshal(offerRequest)
		if err != nil {
			t.Fatalf("Failed to marshal offer request: %v", err)
		}

		req = httptest.NewRequest("POST", "/api/offer", bytes.NewReader(offerBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		apiHandlers.HandleOffer(w, req)

		if w.Code != 204 {
			t.Fatalf("Expected status 204 but got %d", w.Code)
		}

		// Step 3: Retrieve the offer via HTTP
		req = httptest.NewRequest("GET", "/api/offer?token="+token, nil)
		w = httptest.NewRecorder()

		apiHandlers.HandleOffer(w, req)

		if w.Code != 200 {
			t.Fatalf("Expected status 200 but got %d", w.Code)
		}

		var retrievedOffer entities.WebRTCOffer
		err = json.Unmarshal(w.Body.Bytes(), &retrievedOffer)
		if err != nil {
			t.Fatalf("Failed to unmarshal offer response: %v", err)
		}

		if retrievedOffer.Type != "offer" {
			t.Errorf("Expected offer type 'offer' but got %q", retrievedOffer.Type)
		}

		// Step 4: Submit an answer via HTTP
		answerRequest := dto.SubmitAnswerRequest{
			Token: token,
			Answer: &entities.WebRTCAnswer{
				Type: "answer",
				SDP:  "v=0\no=- 987654321 987654321 IN IP4 192.168.1.2\ns=-\nt=0 0\n",
			},
		}

		answerBody, err := json.Marshal(answerRequest)
		if err != nil {
			t.Fatalf("Failed to marshal answer request: %v", err)
		}

		req = httptest.NewRequest("POST", "/api/answer", bytes.NewReader(answerBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		apiHandlers.HandleAnswer(w, req)

		if w.Code != 204 {
			t.Fatalf("Expected status 204 but got %d", w.Code)
		}

		// Step 5: Retrieve the answer via HTTP
		req = httptest.NewRequest("GET", "/api/answer?token="+token, nil)
		w = httptest.NewRecorder()

		apiHandlers.HandleAnswer(w, req)

		if w.Code != 200 {
			t.Fatalf("Expected status 200 but got %d", w.Code)
		}

		var retrievedAnswer entities.WebRTCAnswer
		err = json.Unmarshal(w.Body.Bytes(), &retrievedAnswer)
		if err != nil {
			t.Fatalf("Failed to unmarshal answer response: %v", err)
		}

		if retrievedAnswer.Type != "answer" {
			t.Errorf("Expected answer type 'answer' but got %q", retrievedAnswer.Type)
		}

		// Step 6: Get server info via HTTP
		req = httptest.NewRequest("GET", "/api/info", nil)
		req.Host = "localhost:8080"
		w = httptest.NewRecorder()

		apiHandlers.HandleInfo(w, req)

		if w.Code != 200 {
			t.Fatalf("Expected status 200 but got %d", w.Code)
		}

		var serverInfo entities.ServerInfo
		err = json.Unmarshal(w.Body.Bytes(), &serverInfo)
		if err != nil {
			t.Fatalf("Failed to unmarshal server info response: %v", err)
		}

		if serverInfo.Host != "localhost:8080" {
			t.Errorf("Expected host 'localhost:8080' but got %q", serverInfo.Host)
		}

		if serverInfo.STUNServer != "stun:test.com:19302" {
			t.Errorf("Expected STUN server 'stun:test.com:19302' but got %q", serverInfo.STUNServer)
		}
	})

	t.Run("error scenarios", func(t *testing.T) {
		// Test getting offer for non-existent session
		req := httptest.NewRequest("GET", "/api/offer?token=non-existent", nil)
		w := httptest.NewRecorder()

		apiHandlers.HandleOffer(w, req)

		if w.Code != 404 {
			t.Errorf("Expected status 404 but got %d", w.Code)
		}

		// Test invalid JSON in offer submission
		req = httptest.NewRequest("POST", "/api/offer", bytes.NewReader([]byte("invalid-json")))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()

		apiHandlers.HandleOffer(w, req)

		if w.Code != 400 {
			t.Errorf("Expected status 400 but got %d", w.Code)
		}

		// Test method not allowed
		req = httptest.NewRequest("DELETE", "/api/offer", nil)
		w = httptest.NewRecorder()

		apiHandlers.HandleOffer(w, req)

		if w.Code != 405 {
			t.Errorf("Expected status 405 but got %d", w.Code)
		}

		// Test method not allowed for new token
		req = httptest.NewRequest("GET", "/api/new", nil)
		w = httptest.NewRecorder()

		apiHandlers.HandleNewToken(w, req)

		if w.Code != 405 {
			t.Errorf("Expected status 405 but got %d", w.Code)
		}
	})

	t.Run("concurrent sessions", func(t *testing.T) {
		// Create multiple sessions concurrently
		numSessions := 5
		tokens := make([]string, numSessions)

		for i := 0; i < numSessions; i++ {
			req := httptest.NewRequest("POST", "/api/new", nil)
			w := httptest.NewRecorder()

			apiHandlers.HandleNewToken(w, req)

			if w.Code != 200 {
				t.Fatalf("Failed to create session %d: status %d", i, w.Code)
			}

			var response dto.CreateSessionResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response for session %d: %v", i, err)
			}

			tokens[i] = response.Token
		}

		// Verify all sessions are independent
		for i, token := range tokens {
			for j, otherToken := range tokens {
				if i != j && token == otherToken {
					t.Errorf("Sessions %d and %d have the same token: %s", i, j, token)
				}
			}
		}

		// Submit offers to all sessions
		for i, token := range tokens {
			offerRequest := dto.SubmitOfferRequest{
				Token: token,
				Offer: &entities.WebRTCOffer{
					Type: "offer",
					SDP:  "test-sdp-" + string(rune(i)),
				},
			}

			offerBody, err := json.Marshal(offerRequest)
			if err != nil {
				t.Fatalf("Failed to marshal offer request for session %d: %v", i, err)
			}

			req := httptest.NewRequest("POST", "/api/offer", bytes.NewReader(offerBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			apiHandlers.HandleOffer(w, req)

			if w.Code != 204 {
				t.Fatalf("Failed to submit offer for session %d: status %d", i, w.Code)
			}
		}

		// Verify all offers are stored correctly
		for i, token := range tokens {
			req := httptest.NewRequest("GET", "/api/offer?token="+token, nil)
			w := httptest.NewRecorder()

			apiHandlers.HandleOffer(w, req)

			if w.Code != 200 {
				t.Fatalf("Failed to get offer for session %d: status %d", i, w.Code)
			}

			var offer entities.WebRTCOffer
			err := json.Unmarshal(w.Body.Bytes(), &offer)
			if err != nil {
				t.Fatalf("Failed to unmarshal offer for session %d: %v", i, err)
			}

			expectedSDP := "test-sdp-" + string(rune(i))
			if offer.SDP != expectedSDP {
				t.Errorf("Expected SDP %q for session %d but got %q", expectedSDP, i, offer.SDP)
			}
		}
	})
}