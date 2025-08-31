package resolvers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/internal/services/credential"
	"github.com/weeb-vip/auth/internal/services/mail"
	"github.com/weeb-vip/auth/internal/services/passwordreset"
)

func RequestPasswordReset(
	ctx context.Context,
	authenticationService credential.Credential,
	passwordResetService passwordreset.PasswordReset,
	mailService mail.MailService,
	cfg *config.Config,
	username string,
	email string,
) (bool, error) {
	foundCredential, err := authenticationService.GetCredentials(ctx, username)
	if err != nil {
		return false, err
	}

	reset, err := passwordResetService.PasswordResetRequest(ctx, foundCredential.ID)
	if err != nil {
		return false, err
	}

	resetURL, err := url.Parse(cfg.APPConfig.PasswordResetBaseURL)
	if err != nil {
		return false, fmt.Errorf("invalid password reset base URL: %w", err)
	}
	
	query := resetURL.Query()
	query.Set("token", reset.OTT)
	query.Set("username", username)
	resetURL.RawQuery = query.Encode()

	err = mailService.SendMail(ctx, []string{email}, "Password Reset", "reset-password.mjml", map[string]string{
		"token_url": resetURL.String(),
		"name":      username,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
