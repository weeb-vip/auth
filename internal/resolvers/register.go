package resolvers

import (
	"context"
	"fmt"
	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/internal/services/mail"
	"log"
	"net/url"

	"github.com/weeb-vip/auth/graph/model"
	"github.com/weeb-vip/auth/internal/services/validation_token"

	"github.com/weeb-vip/auth/internal/services/credential"
)

func Register( // nolint
	ctx context.Context,
	cfg *config.Config,
	authenticationService credential.Credential,
	validatonToken validation_token.ValidationToken,
	mailService mail.MailService,
	username string,
	password string,
) (*model.RegisterResult, error) {
	credentials, err := authenticationService.Register(ctx, username, password)
	if err != nil {
		res, err := handleError(ctx, "null", err)
		if res != nil {
			return nil, err
		}

		return nil, err
	}

	token, err := validatonToken.GenerateToken(credentials.ID)
	if err != nil {
		return nil, err
	}

	resetURL, err := url.Parse(cfg.APPConfig.VerificationBaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid password reset base URL: %w", err)
	}

	query := resetURL.Query()
	query.Set("token", token)
	query.Set("email", username)
	resetURL.RawQuery = query.Encode()

	err = mailService.SendMail(ctx, []string{username}, "Password Reset", "verification.mjml", map[string]string{
		"token_url": resetURL.String(),
		"name":      username,
	})

	log.Println(token)

	return &model.RegisterResult{
		ID: credentials.UserID,
	}, nil
}
