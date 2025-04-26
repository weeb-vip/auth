package session

import (
	"context"

	"github.com/weeb-vip/auth/internal/services/session/models"
	"github.com/weeb-vip/auth/internal/services/session/repositories"
)

const (
	ErrSessionNotFound = "session not found"
)

type sessionService struct {
	sessionRepository repositories.SessionsRepository
}

func NewSessionService() Session {
	sessionRepository := repositories.NewSessionsRepository()

	return &sessionService{
		sessionRepository: sessionRepository,
	}
}

func (service *sessionService) CreateSession(
	ctx context.Context,
	username string,
) (*models.Session, error) {
	return service.sessionRepository.CreateSession(ctx, username)
}

func (service *sessionService) GetSession(
	ctx context.Context,
	token string,
) (*models.Session, error) {
	session, err := service.sessionRepository.GetSession(ctx, token)
	if err != nil {
		return nil, err
	}

	if session == nil {
		return nil, &Error{
			Code:    ErrSessionNotFound,
			Message: "session not found",
		}
	}

	return session, nil
}

func (service *sessionService) DeleteSession(
	ctx context.Context,
	token string,
) error {
	return service.sessionRepository.DeleteSession(ctx, token)
}
