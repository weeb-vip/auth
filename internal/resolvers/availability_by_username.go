package resolvers

import (
	"context"

	"github.com/weeb-vip/auth/internal/services/credential"
)

func AvailabilityByUsername(
	ctx context.Context,
	authenticationService credential.Credential,
	username string,
) (bool, error) {
	cred, err := authenticationService.GetCredentials(ctx, username)
	if err != nil {
		return false, err
	}

	if cred != nil {
		return false, nil
	}

	return true, nil
}
