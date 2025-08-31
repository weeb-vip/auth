package resolvers

import (
	"context"
	"fmt"
	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/http/handlers/logger"
	"github.com/weeb-vip/auth/http/handlers/requestinfo"
	"github.com/weeb-vip/auth/internal/services/credential"
	"github.com/weeb-vip/auth/internal/services/mail"
	"github.com/weeb-vip/auth/internal/services/validation_token"
	"net/url"
)

func EmailVerification( // nolint
	ctx context.Context,
	credentialService credential.Credential,
) (bool, error) {
	log := logger.FromContext(ctx)
	req := requestinfo.FromContext(ctx)

	userID := req.UserID
	if userID == nil {
		log.Error("User ID is missing")
		return false, nil
	}

	err := credentialService.ActivateCredentials(ctx, *userID)
	if err != nil {
		res, err := handleError(ctx, "false", err)
		if res != nil {
			return false, err
		}

		return false, err
	}

	return true, nil
}

func ResendVerificationEmail( // nolint
	ctx context.Context,
	credentialService credential.Credential,
	validatonToken validation_token.ValidationToken,
	mailService mail.MailService,
	cfg *config.Config,
	username string,
) (bool, error) {
	log := logger.FromContext(ctx)

	credential, err := credentialService.GetCredentials(ctx, username)
	if err != nil {
		res, err := handleError(ctx, "false", err)
		if res != nil {
			return false, err
		}

		return false, err
	}

	if credential == nil {
		log.Error("Credential not found")
		return false, nil
	}

	if credential.Active {
		log.Error("Credential is already active")
		return false, nil
	}

	token, err := validatonToken.GenerateToken(credential.ID)
	if err != nil {
		return false, err
	}

	resetURL, err := url.Parse(cfg.APPConfig.VerificationBaseURL)
	if err != nil {
		return false, fmt.Errorf("invalid password reset base URL: %w", err)
	}

	query := resetURL.Query()
	query.Set("token", token)
	query.Set("email", username)
	resetURL.RawQuery = query.Encode()

	err = mailService.SendMail(ctx, []string{username}, "Email Verification", "verification.mjml", map[string]string{
		"token_url": resetURL.String(),
		"name":      username,
	})

	return true, nil
}
