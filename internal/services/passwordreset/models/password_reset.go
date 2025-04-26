package models

import (
	"github.com/weeb-vip/auth/internal/db"
)

type PasswordReset struct {
	db.BaseModel
	CredentialID string `json:"credentialId"`
	OTT          string `json:"ott"`
}
