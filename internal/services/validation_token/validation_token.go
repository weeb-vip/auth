package validation_token // nolint

import (
	"fmt"
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

func (v validationToken) ValidateToken(token string) (*string, error) {
	// get claims from token
	claims, err := v.Tokenizer.GetClaims(token)
	if err != nil {
		return nil, err
	}

	if claims.Purpose == nil || *claims.Purpose != "EMAIL_VERIFICATION" {
		return nil, fmt.Errorf("invalid token purpose")
	}

	return claims.Subject, nil

}
