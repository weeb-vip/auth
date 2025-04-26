package refresh_token_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/internal/services/refresh_token"
)

func TestNewRefreshTokenService(t *testing.T) {
	t.Parallel()
	t.Run("should return a new refresh token service", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)
		cfg, _ := config.LoadConfig()
		refreshTokenService := refresh_token.NewRefreshTokenService(cfg.RefreshTokenConfig)

		a.NotNil(refreshTokenService)
	})

	t.Run("should get token", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)
		cfg, _ := config.LoadConfig()
		refreshTokenService := refresh_token.NewRefreshTokenService(cfg.RefreshTokenConfig)

		_, err := refreshTokenService.GetToken("token")
		a.NoError(err)
	})

	t.Run("should delete token", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)
		cfg, _ := config.LoadConfig()
		refreshTokenService := refresh_token.NewRefreshTokenService(cfg.RefreshTokenConfig)

		err := refreshTokenService.DeleteToken("token")
		a.NoError(err)
	})

	t.Run("should create token", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)
		cfg, _ := config.LoadConfig()
		refreshTokenService := refresh_token.NewRefreshTokenService(cfg.RefreshTokenConfig)

		token, err := refreshTokenService.CreateToken("userid")
		a.NoError(err)
		a.NotEmpty(token)
	})

	t.Run("Should validate token", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)
		cfg, _ := config.LoadConfig()
		refreshTokenService := refresh_token.NewRefreshTokenService(cfg.RefreshTokenConfig)

		token, err := refreshTokenService.CreateToken("userid")
		a.NoError(err)
		a.NotEmpty(token)
		a.True(refreshTokenService.ValidateToken(token.Token))
	})
}
