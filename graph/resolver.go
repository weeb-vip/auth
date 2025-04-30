package graph

import (
	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/internal/jwt"
	"github.com/weeb-vip/auth/internal/services/credential"
	"github.com/weeb-vip/auth/internal/services/passwordreset"
	"github.com/weeb-vip/auth/internal/services/refresh_token"
	"github.com/weeb-vip/auth/internal/services/session"
	"github.com/weeb-vip/auth/internal/services/validation_token"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CredentialService    credential.Credential
	PasswordResetService passwordreset.PasswordReset
	SessionService       session.Session
	JwtTokenizer         jwt.Tokenizer
	Config               config.Config
	RefreshTokenService  refresh_token.RefreshToken
	ValidationToken      validation_token.ValidationToken
}
