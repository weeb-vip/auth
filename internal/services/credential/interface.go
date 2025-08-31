package credential

import (
	"context"

	"github.com/weeb-vip/auth/internal/services/credential/models"
)

type Credential interface {
	Register(ctx context.Context, username string, password string) (*models.Credential, error)
	SignIn(ctx context.Context, username string, password string) (*models.Credential, error)
	GetCredentials(ctx context.Context, username string) (*models.Credential, error)
	UpdatePassword(ctx context.Context, username string, newPassword string) error
	ActivateCredentials(ctx context.Context, identifier string) error
	GetCredentialsByIdentifier(ctx context.Context, identifier string) (*models.Credential, error)
}
