package credential

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     int = 4  // The minimum allowable cost as passed in to GenerateFromPassword.
	MaxCost     int = 14 // The maximum allowable cost as passed in to GenerateFromPassword.
	DefaultCost int = 10 // The cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword.
)

func (service *credentialService) HashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), MaxCost)
	if err != nil {
		return "", err
	}

	return string(pass), nil
}

func (service *credentialService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
