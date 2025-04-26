package passwordreset

import (
	"context"

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
