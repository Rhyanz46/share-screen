package http

import (
	"encoding/json"
	"log"
	"net/http"

	"share-screen/pkg/domain/interfaces"
	"share-screen/pkg/usecase/dto"
	"share-screen/pkg/usecase/usecases"
)

// APIHandlers contains handlers for API endpoints
type APIHandlers struct {
	sessionUseCase    interfaces.SessionUseCase
	serverInfoUseCase interfaces.ServerInfoUseCase
}

// NewAPIHandlers creates a new API handlers instance
func NewAPIHandlers(sessionUseCase interfaces.SessionUseCase, serverInfoUseCase interfaces.ServerInfoUseCase) *APIHandlers {
	return &APIHandlers{
		sessionUseCase:    sessionUseCase,
		serverInfoUseCase: serverInfoUseCase,
	}
}

// HandleNewToken creates a new session token
func (h *APIHandlers) HandleNewToken(w http.ResponseWriter, r *http.Request) {
	log.Printf("📞 API: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	response, err := h.sessionUseCase.CreateSession()
	if err != nil {
		log.Printf("❌ Error creating session: %v", err)
		http.Error(w, "failed to generate token", 500)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding token response: %v", err)
		http.Error(w, "internal server error", 500)
	}
}

// HandleOffer handles WebRTC offer operations (POST to store, GET to retrieve)
func (h *APIHandlers) HandleOffer(w http.ResponseWriter, r *http.Request) {
	log.Printf("📞 API: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	switch r.Method {
	case http.MethodPost:
		h.handleSubmitOffer(w, r)
	case http.MethodGet:
		h.handleGetOffer(w, r)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func (h *APIHandlers) handleSubmitOffer(w http.ResponseWriter, r *http.Request) {
	var request dto.SubmitOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("❌ Invalid offer payload: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}

	log.Printf("🔴 Sender posting offer for token: %s...", request.Token[:8])

	if err := h.sessionUseCase.SubmitOffer(&request); err != nil {
		log.Printf("❌ Error submitting offer: %v", err)
		h.handleUseCaseError(w, err)
		return
	}

	w.WriteHeader(204)
}

func (h *APIHandlers) handleGetOffer(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	log.Printf("🔵 Viewer requesting offer for token: %s...", token[:8])

	request := &dto.GetOfferRequest{Token: token}
	response, err := h.sessionUseCase.GetOffer(request)
	if err != nil {
		h.handleUseCaseError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response.Offer); err != nil {
		log.Printf("Error encoding offer response: %v", err)
		http.Error(w, "internal server error", 500)
	}
}

// HandleAnswer handles WebRTC answer operations (POST to store, GET to retrieve)
func (h *APIHandlers) HandleAnswer(w http.ResponseWriter, r *http.Request) {
	log.Printf("📞 API: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	switch r.Method {
	case http.MethodPost:
		h.handleSubmitAnswer(w, r)
	case http.MethodGet:
		h.handleGetAnswer(w, r)
	default:
		http.Error(w, "method not allowed", 405)
	}
}

func (h *APIHandlers) handleSubmitAnswer(w http.ResponseWriter, r *http.Request) {
	var request dto.SubmitAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("❌ Invalid answer payload: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}

	log.Printf("🔵 Viewer posting answer for token: %s...", request.Token[:8])

	if err := h.sessionUseCase.SubmitAnswer(&request); err != nil {
		log.Printf("❌ Error submitting answer: %v", err)
		h.handleUseCaseError(w, err)
		return
	}

	w.WriteHeader(204)
}

func (h *APIHandlers) handleGetAnswer(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	log.Printf("🔴 Sender requesting answer for token: %s...", token[:8])

	request := &dto.GetAnswerRequest{Token: token}
	response, err := h.sessionUseCase.GetAnswer(request)
	if err != nil {
		h.handleUseCaseError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response.Answer); err != nil {
		log.Printf("Error encoding answer response: %v", err)
		http.Error(w, "internal server error", 500)
	}
}

// HandleInfo provides server information including LAN IP
func (h *APIHandlers) HandleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	serverInfo, err := h.serverInfoUseCase.GetServerInfo(r.Host)
	if err != nil {
		log.Printf("Error getting server info: %v", err)
		http.Error(w, "internal server error", 500)
		return
	}

	if err := json.NewEncoder(w).Encode(serverInfo); err != nil {
		log.Printf("Error encoding info response: %v", err)
		http.Error(w, "internal server error", 500)
	}
}

// handleUseCaseError converts use case errors to appropriate HTTP responses
func (h *APIHandlers) handleUseCaseError(w http.ResponseWriter, err error) {
	switch err {
	case usecases.ErrSessionNotFound:
		http.Error(w, "session not found", 404)
	case usecases.ErrSessionExpired:
		http.Error(w, "session expired", 410)
	case usecases.ErrInvalidOffer, usecases.ErrInvalidAnswer:
		http.Error(w, err.Error(), 400)
	case usecases.ErrOfferNotFound:
		http.Error(w, "offer not found", 404)
	case usecases.ErrAnswerNotFound:
		http.Error(w, "answer not found", 404)
	case usecases.ErrAnswerAlreadyExists:
		http.Error(w, "answer already exists", 409)
	case usecases.ErrSessionNotReady:
		http.Error(w, "session not ready", 400)
	default:
		log.Printf("Unexpected error: %v", err)
		http.Error(w, "internal server error", 500)
	}
}
