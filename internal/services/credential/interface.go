package credential

import (
	"context"

	"github.com/weeb-vip/auth/internal/services/credential/models"
)

type Credential interface {
	Register(ctx context.Context, username string, password string) (*models.Credential, error)
	SignIn(ctx context.Context, username string, password string) (*models.Credential, error)
	GetCredentials(ctx context.Context, username string) (*models.Credential, error)
}
