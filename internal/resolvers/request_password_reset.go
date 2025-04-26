package resolvers

import (
	"context"

	"github.com/weeb-vip/auth/internal/services/passwordreset"

	"github.com/weeb-vip/auth/internal/services/credential"
)

func RequestPasswordReset(
	ctx context.Context,
	authenticationService credential.Credential,
	passwordResetService passwordreset.PasswordReset,
	username string,
) (bool, error) {
	foundCredential, err := authenticationService.GetCredentials(ctx, username)
	if err != nil {
		return false, err
	}

	_, err = passwordResetService.PasswordResetRequest(ctx, foundCredential.ID)

	if err != nil {
		return false, err
	}

	return true, nil
}
