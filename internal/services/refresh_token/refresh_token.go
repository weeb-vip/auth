package refresh_token //nolint

import (
	"time"

	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/internal/services/refresh_token/models"
	"github.com/weeb-vip/auth/internal/services/refresh_token/repositories"
	"github.com/weeb-vip/auth/internal/ulid"
)

type refreshTokenService struct {
	refreshTokenRepository repositories.RefreshTokenRepository
	config                 config.RefreshTokenConfig
}

func NewRefreshTokenService(config config.RefreshTokenConfig) RefreshToken {
	refreshTokenRepository := repositories.GetRefreshTokenRepository()

	return &refreshTokenService{
		config:                 config,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (service *refreshTokenService) GetTokenByUserID(userID string) (*models.RefreshToken, error) {
	return service.refreshTokenRepository.GetRefreshTokenByUserID(userID)
}

func (service *refreshTokenService) GetToken(token string) (*models.RefreshToken, error) {
	return service.refreshTokenRepository.GetRefreshToken(token)
}

func (service *refreshTokenService) CreateToken(userID string) (*models.RefreshToken, error) {
	ttl := time.Duration(service.config.TokenTTL) * time.Hour
	expiry := time.Now().Add(ttl).Unix()

	return service.refreshTokenRepository.AddRefreshToken(userID, ulid.New("refresh_token"), expiry)
}

func (service *refreshTokenService) DeleteToken(userID string) error {
	return service.refreshTokenRepository.DeleteRefreshToken(userID)
}

func (service *refreshTokenService) ValidateToken(token string) (bool, error) {
	refreshToken, err := service.GetToken(token)
	if err != nil {
		return false, err
	}

	if refreshToken.Expiry < time.Now().Unix() {
		return false, nil
	}

	return true, nil
}
