package usecases

import (
	"errors"
	"log"
	"time"

	"share-screen/pkg/domain/entities"
	"share-screen/pkg/domain/interfaces"
	"share-screen/pkg/usecase/dto"
)

var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrSessionExpired      = errors.New("session expired")
	ErrInvalidOffer        = errors.New("invalid offer")
	ErrInvalidAnswer       = errors.New("invalid answer")
	ErrOfferNotFound       = errors.New("offer not found")
	ErrAnswerNotFound      = errors.New("answer not found")
	ErrAnswerAlreadyExists = errors.New("answer already exists")
	ErrSessionNotReady     = errors.New("session not ready for answer")
)

// SessionUseCase implements the session use case interface
type SessionUseCase struct {
	sessionRepo interfaces.SessionRepository
	tokenExpiry time.Duration
}

// NewSessionUseCase creates a new session use case
func NewSessionUseCase(sessionRepo interfaces.SessionRepository, tokenExpiry time.Duration) *SessionUseCase {
	return &SessionUseCase{
		sessionRepo: sessionRepo,
		tokenExpiry: tokenExpiry,
	}
}

// CreateSession creates a new screen sharing session
func (uc *SessionUseCase) CreateSession() (*dto.CreateSessionResponse, error) {
	session, err := uc.sessionRepo.CreateSession(uc.tokenExpiry)
	if err != nil {
		log.Printf("❌ Error creating session: %v", err)
		return nil, err
	}

	log.Printf("🚀 Sender session started with token: %s...", session.Token[:8])

	return &dto.CreateSessionResponse{
		Token: session.Token,
	}, nil
}

// SubmitOffer submits a WebRTC offer for a session
func (uc *SessionUseCase) SubmitOffer(request *dto.SubmitOfferRequest) error {
	if request.Offer == nil || !request.Offer.IsValid() {
		return ErrInvalidOffer
	}

	session, err := uc.sessionRepo.GetSession(request.Token)
	if err != nil {
		return ErrSessionNotFound
	}

	if session.IsExpired() {
		return ErrSessionExpired
	}

	if !session.CanAcceptOffer() {
		return errors.New("session cannot accept offer")
	}

	session.Offer = request.Offer
	session.Status = entities.SessionStatusActive

	if err := uc.sessionRepo.UpdateSession(session); err != nil {
		log.Printf("❌ Error updating session with offer: %v", err)
		return err
	}

	log.Printf("📤 Offer created for token: %s (type: %s)", request.Token[:8]+"...", request.Offer.Type)
	return nil
}

// GetOffer retrieves a WebRTC offer for a session
func (uc *SessionUseCase) GetOffer(request *dto.GetOfferRequest) (*dto.GetOfferResponse, error) {
	session, err := uc.sessionRepo.GetSession(request.Token)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.IsExpired() {
		return nil, ErrSessionExpired
	}

	if session.Offer == nil {
		log.Printf("❌ Offer not found for token: %s", request.Token[:8]+"...")
		return nil, ErrOfferNotFound
	}

	log.Printf("📥 Offer retrieved for token: %s", request.Token[:8]+"...")
	return &dto.GetOfferResponse{
		Offer: session.Offer,
	}, nil
}

// SubmitAnswer submits a WebRTC answer for a session
func (uc *SessionUseCase) SubmitAnswer(request *dto.SubmitAnswerRequest) error {
	if request.Answer == nil || !request.Answer.IsValid() {
		return ErrInvalidAnswer
	}

	session, err := uc.sessionRepo.GetSession(request.Token)
	if err != nil {
		return ErrSessionNotFound
	}

	if session.IsExpired() {
		return ErrSessionExpired
	}

	if !session.CanAcceptAnswer() {
		if session.Answer != nil {
			log.Printf("⚠️  Answer already exists for token: %s", request.Token[:8]+"...")
			return ErrAnswerAlreadyExists
		}
		return ErrSessionNotReady
	}

	session.Answer = request.Answer

	if err := uc.sessionRepo.UpdateSession(session); err != nil {
		log.Printf("❌ Error updating session with answer: %v", err)
		return err
	}

	log.Printf("📤 Answer created for token: %s (type: %s)", request.Token[:8]+"...", request.Answer.Type)
	log.Printf("🎯 WebRTC handshake completed for token: %s", request.Token[:8]+"...")
	return nil
}

// GetAnswer retrieves a WebRTC answer for a session
func (uc *SessionUseCase) GetAnswer(request *dto.GetAnswerRequest) (*dto.GetAnswerResponse, error) {
	session, err := uc.sessionRepo.GetSession(request.Token)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.IsExpired() {
		return nil, ErrSessionExpired
	}

	if session.Answer == nil {
		log.Printf("❌ Answer not ready for token: %s", request.Token[:8]+"...")
		return nil, ErrAnswerNotFound
	}

	log.Printf("📥 Answer retrieved for token: %s", request.Token[:8]+"...")
	return &dto.GetAnswerResponse{
		Answer: session.Answer,
	}, nil
}
