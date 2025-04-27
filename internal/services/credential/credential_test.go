package credential_test

import (
	"context"
	"github.com/weeb-vip/auth/mocks"
	"go.uber.org/mock/gomock"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/weeb-vip/auth/internal/services/credential"
)

func TestCredentialService_Register(t *testing.T) {
	t.Parallel()
	t.Run("Test Register", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)
		ctrl := gomock.NewController(t)

		userClient := mocks.NewMockUserClientInterface(ctrl)
		credentialService := credential.NewCredentialService(userClient)

		_, err := credentialService.Register(context.TODO(), "username", "password")
		a.NoError(err)
	})
	t.Run("Test Register 2 times - idempotence", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)
		ctrl := gomock.NewController(t)

		userClient := mocks.NewMockUserClientInterface(ctrl)
		credentialService := credential.NewCredentialService(userClient)

		_, err := credentialService.Register(context.TODO(), "username2", "password")
		a.NoError(err)
		_, err = credentialService.Register(context.TODO(), "username2", "password")
		a.NoError(err)
	})
}

func TestCredentialService_SignIn(t *testing.T) {
	t.Parallel()
	t.Run("Test SignIn", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)
		ctrl := gomock.NewController(t)

		userClient := mocks.NewMockUserClientInterface(ctrl)
		credentialService := credential.NewCredentialService(userClient)

		_, err := credentialService.Register(context.TODO(), "username", "password")
		a.NoError(err)

		a.NotNil(credentialService.SignIn(context.TODO(), "username", "password"))
		a.Nil(credentialService.SignIn(context.TODO(), "username", "password2"))
	})
}
