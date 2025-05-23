package resolvers

import (
	"context"

	"github.com/weeb-vip/auth/graph/model"
	"github.com/weeb-vip/auth/internal/jwt"
	"github.com/weeb-vip/auth/internal/services/refresh_token"

	"github.com/weeb-vip/auth/internal/services/credential"
	"github.com/weeb-vip/auth/internal/services/session"
	SessionModels "github.com/weeb-vip/auth/internal/services/session/models"
	"github.com/weeb-vip/auth/internal/ulid"
)

func CreateSession( // nolint
	ctx context.Context,
	credentialService credential.Credential,
	sessionService session.Session,
	refreshTokenService refresh_token.RefreshToken,
	jwtTokenizer jwt.Tokenizer,
	input *model.LoginInput,
) (*model.SigninResult, error) {
	createdSession, err := createSession(ctx, input, sessionService, credentialService)

	if err != nil {
		_, err := handleError(ctx, "null", err)

		return nil, err
	}

	subject := createdSession.UserID

	refreshToken, err := refreshTokenService.CreateToken(subject)
	if err != nil {
		return nil, err
	}

	token, err := jwtTokenizer.Tokenize(jwt.Claims{
		Subject:      &subject,
		TTL:          nil,
		Purpose:      nil,
		RefreshToken: &refreshToken.Token,
	})

	if err != nil {
		return nil, err
	}

	return &model.SigninResult{
		ID: createdSession.UserID,
		Credentials: &model.Credentials{
			Token:        token,
			RefreshToken: refreshToken.Token,
		},
	}, nil
}

func createSession(ctx context.Context,
	input *model.LoginInput,
	sessionService session.Session,
	credentialService credential.Credential,
) (*SessionModels.Session, error) {
	if input == nil {
		return sessionService.CreateSession(ctx, ulid.New("guest"))
	}

	return createUserSession(ctx, sessionService, credentialService, *input)
}

func createUserSession(
	ctx context.Context,
	sessionService session.Session,
	credentialService credential.Credential,
	input model.LoginInput,
) (*SessionModels.Session, error) {
	result, err := credentialService.SignIn(ctx, input.Username, input.Password)

	if err != nil {
		return nil, err
	}

	return sessionService.CreateSession(ctx, result.UserID)
}
