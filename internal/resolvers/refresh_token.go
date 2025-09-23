package resolvers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/graph/model"
	"github.com/weeb-vip/auth/http/handlers/responsecontext"
	"github.com/weeb-vip/auth/internal/jwt"
	"github.com/weeb-vip/auth/internal/services/refresh_token"
	"github.com/weeb-vip/auth/internal/services/session"
)

func RefreshToken(ctx context.Context, sessionService session.Session, refreshTokenService refresh_token.RefreshToken, jwtTokenizer jwt.Tokenizer, config *config.Config, token string) (*model.SigninResult, error) {

	refreshToken, err := refreshTokenService.GetToken(token)
	if err != nil {
		return nil, err
	}
	if refreshToken == nil {
		return nil, fmt.Errorf("refresh token not found")
	}

	session, err := sessionService.CreateSession(ctx, refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	subject := session.UserID

	refreshToken, err = refreshTokenService.CreateToken(subject)
	if err != nil {
		return nil, err
	}

	token, err = jwtTokenizer.Tokenize(jwt.Claims{
		Subject:      &subject,
		TTL:          nil,
		Purpose:      nil,
		RefreshToken: &refreshToken.Token,
	})

	if err != nil {
		return nil, err
	}

	// Set access token as HTTP-only cookie
	responseWriter := responsecontext.FromContext(ctx)

	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		Domain:   config.APPConfig.CookieDomain,
		HttpOnly: true,
		Secure:   true, // Allow non-HTTPS for development
		SameSite: http.SameSiteNoneMode,
		MaxAge:   int((time.Minute * 15).Seconds()), // 15 minutes
	}

	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Token,
		Path:     "/",
		Domain:   config.APPConfig.CookieDomain,
		HttpOnly: true,
		Secure:   true, // Allow non-HTTPS for development
		SameSite: http.SameSiteNoneMode,
		MaxAge:   int((time.Hour * 24 * 7).Seconds()), // 7 days
	}

	http.SetCookie(responseWriter, accessTokenCookie)
	http.SetCookie(responseWriter, refreshTokenCookie)

	return &model.SigninResult{
		ID: session.ID,
		Credentials: &model.Credentials{
			Token:        &token,
			RefreshToken: &refreshToken.Token,
		},
	}, nil
}
