package passwordreset

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/weeb-vip/auth/internal/services/passwordreset/models"
	"github.com/weeb-vip/auth/internal/services/passwordreset/repositories"
)

type passwordResetService struct {
	passwordResetRepository repositories.PasswordResetRepository
}

func NewPasswordResetService() PasswordReset {
	passwordResetRepository := repositories.NewPasswordResetRepository()

	return &passwordResetService{
		passwordResetRepository: passwordResetRepository,
	}
}

func (service *passwordResetService) PasswordResetRequest(
	ctx context.Context,
	credentialID string,
) (*models.PasswordReset, error) {
	ott := uuid.New().String()

	return service.passwordResetRepository.AddOTT(credentialID, ott)
}

func (service *passwordResetService) ValidateAndConsumeToken(
	ctx context.Context,
	username string,
	token string,
) error {
	passwordReset, err := service.passwordResetRepository.GetOTTByToken(token)
	if err != nil {
		return err
	}

	if passwordReset.OTT == "" {
		return errors.New("invalid token")
	}

	err = service.passwordResetRepository.DeleteOTTByToken(token)
	if err != nil {
		return err
	}

	return nil
}
