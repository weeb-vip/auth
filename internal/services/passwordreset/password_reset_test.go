package passwordreset

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weeb-vip/auth/internal/services/credential"
)

func TestPasswordResetService_PasswordResetRequest(t *testing.T) {
	t.Run("Test PasswordResetRequest", func(t *testing.T) {
		a := assert.New(t)

		credentialService := credential.NewCredentialService()

		cred, err := credentialService.Register(context.TODO(), "username", "password")

		passwordResetService := NewPasswordResetService()
		_, err = passwordResetService.PasswordResetRequest(context.TODO(), cred.ID)
		a.NoError(err)
	})
}
