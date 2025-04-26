package session

import (
	"context"

	"github.com/weeb-vip/auth/internal/services/session/models"
)

type Session interface {
	CreateSession(ctx context.Context, userID string) (*models.Session, error)
}
