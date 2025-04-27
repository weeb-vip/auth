package validation_token_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/weeb-vip/auth/internal/services/validation_token"
	"github.com/weeb-vip/auth/mocks"
)

func TestValidationToken_GenerateToken(t *testing.T) {
	t.Run("generate token", func(t *testing.T) {
		a := assert.New(t)
		ctrl := gomock.NewController(t)
		tokenizer := mocks.NewMockTokenizer(ctrl)

		tokenizer.EXPECT().Tokenize(gomock.Any()).Return("jwt", nil)

		validationToken := validation_token.NewValidationTokenService(tokenizer)

		jwt, err := validationToken.GenerateToken("identifier")
		a.NoError(err)
		a.NotNil(jwt)
	})
}
