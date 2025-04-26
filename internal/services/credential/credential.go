package credential

import (
	"context"
	"github.com/weeb-vip/auth/internal/services/user_client"

	"github.com/weeb-vip/auth/internal/services/credential/models"
	"github.com/weeb-vip/auth/internal/services/credential/repositories"
	"github.com/weeb-vip/auth/internal/ulid"
)

const (
	ErrPasswordMismatch = "ERROR_PASSWORD_MISMATCH"
)

type credentialService struct {
	credentialsRepository repositories.CredentialsRepository
	userClient            user_client.UserClientInterface
}

func NewCredentialService(userClient user_client.UserClientInterface) Credential {
	credentialRepository := repositories.GetCredentialsRepository()

	return &credentialService{
		credentialsRepository: credentialRepository,
		userClient:            userClient,
	}
}

func (service *credentialService) GetCredentials(ctx context.Context, username string) (*models.Credential, error) {
	return service.credentialsRepository.GetCredentials(username)
}

func (service *credentialService) Register(
	ctx context.Context,
	username string,
	password string,
) (*models.Credential, error) {
	hashedPassword, err := service.HashPassword(password)
	if err != nil {
		return nil, &Error{
			Code:    CredentialErrorInternalError,
			Message: err.Error(),
		}
	}

	userID := ulid.New("user")

	result, err := service.credentialsRepository.AddCredentials(
		username,
		userID,
		hashedPassword,
		models.PasswordCredential,
	)

	if err != nil {
		return nil, &Error{
			Code:    CredentialErrorInternalError,
			Message: err.Error(),
		}
	}

	return result, nil
}

func (service *credentialService) SignIn( //nolint
	ctx context.Context,
	username string,
	password string,
) (*models.Credential, error) {
	credentials, err := service.credentialsRepository.GetCredentials(username) // nolint
	if err != nil {
		return nil, &Error{
			Code:    CredentialErrorInternalError,
			Message: "database error",
		}
	}

	if credentials == nil {
		return nil, &Error{
			Code:    CredentialErrorInvalidCredentials,
			Message: "invalid credentials",
		}
	}

	if service.VerifyPassword(password, credentials.Value) {
		return credentials, nil
	}

	return nil, &Error{
		Code:    CredentialErrorInvalidCredentials,
		Message: "invalid credentials",
	}
}
