package validation_token // nolint

import (
	"github.com/weeb-vip/auth/internal/jwt"
)

type validationToken struct {
	Tokenizer jwt.Tokenizer
}

func NewValidationTokenService(tokenizer jwt.Tokenizer) ValidationToken {
	return validationToken{
		tokenizer,
	}
}

func (v validationToken) GenerateToken(identifier string) (string, error) {
	purpose := "EMAIL_VERIFICATION"

	return v.Tokenizer.Tokenize(jwt.Claims{
		Subject: &identifier,
		TTL:     nil,
		Purpose: &purpose,
	})
}
