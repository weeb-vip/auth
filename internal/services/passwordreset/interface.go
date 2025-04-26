package passwordreset

import (
	"context"

	"github.com/weeb-vip/auth/internal/services/passwordreset/models"
)

type PasswordReset interface {
	PasswordResetRequest(ctx context.Context, credentialID string) (*models.PasswordReset, error)
}
