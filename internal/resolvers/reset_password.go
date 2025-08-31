package resolvers

import (
	"context"

	"github.com/weeb-vip/auth/internal/services/credential"
	"github.com/weeb-vip/auth/internal/services/passwordreset"
)

func ResetPassword(
	ctx context.Context,
	credentialService credential.Credential,
	passwordResetService passwordreset.PasswordReset,
	token string,
	username string,
	newPassword string,
) (bool, error) {
	err := passwordResetService.ValidateAndConsumeToken(ctx, username, token)
	if err != nil {
		return false, err
	}

	err = credentialService.UpdatePassword(ctx, username, newPassword)
	if err != nil {
		return false, err
	}

	return true, nil
}