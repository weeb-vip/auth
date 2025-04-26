package refresh_token //nolint

import "github.com/weeb-vip/auth/internal/services/refresh_token/models"

type RefreshToken interface {
	GetToken(token string) (*models.RefreshToken, error)
	GetTokenByUserID(userID string) (*models.RefreshToken, error)
	CreateToken(userID string) (*models.RefreshToken, error)
	DeleteToken(userID string) error
	ValidateToken(token string) (bool, error)
}
