package resolvers

import (
	"context"
	"github.com/weeb-vip/auth/internal/services/credential"
	"github.com/weeb-vip/auth/internal/services/validation_token"
)

func EmailVerification( // nolint
	ctx context.Context,
	validatonToken validation_token.ValidationToken,
	credentialService credential.Credential,
	token string,
) (bool, error) {
	identifier, err := validatonToken.ValidateToken(token)
	if err != nil {
		res, err := handleError(ctx, "false", err)
		if res != nil {
			return false, err
		}

		return false, err
	}

	err = credentialService.ActivateCredentials(ctx, *identifier)
	if err != nil {
		res, err := handleError(ctx, "false", err)
		if res != nil {
			return false, err
		}

		return false, err
	}

	return true, nil
}
