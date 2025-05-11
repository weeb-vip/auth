package resolvers

import (
	"context"
	"fmt"
	"github.com/weeb-vip/auth/graph/model"
	"github.com/weeb-vip/auth/internal/jwt"
	"github.com/weeb-vip/auth/internal/services/refresh_token"
	"github.com/weeb-vip/auth/internal/services/session"
)

func RefreshToken(ctx context.Context, sessionService session.Session, refreshTokenService refresh_token.RefreshToken, jwtTokenizer jwt.Tokenizer, token string) (*model.SigninResult, error) {

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
	return &model.SigninResult{
		ID: session.ID,
		Credentials: &model.Credentials{
			Token:        token,
			RefreshToken: refreshToken.Token,
		},
	}, nil
}
