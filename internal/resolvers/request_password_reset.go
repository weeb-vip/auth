package resolvers

import (
	"context"
	"github.com/weeb-vip/auth/internal/services/mail"
	"github.com/weeb-vip/auth/internal/services/passwordreset"

	"github.com/weeb-vip/auth/internal/services/credential"
)

func RequestPasswordReset(
	ctx context.Context,
	authenticationService credential.Credential,
	passwordResetService passwordreset.PasswordReset,
	mailService mail.MailService,
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

	token_url := reset.OTT

	err = mailService.SendMail(ctx, []string{email}, "Password Reset", "reset-password.mjml", map[string]string{
		"token_url": token_url,
		"name":      username,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
