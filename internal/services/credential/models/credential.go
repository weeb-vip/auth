package models

import (
	"github.com/weeb-vip/auth/internal/db"
)

const PasswordCredential CredentialTypes = "password"
const TokenCredential CredentialTypes = "token" // nolint

type CredentialTypes string

type Credential struct {
	db.BaseModel
	Username string          `json:"username"`
	UserID   string          `column:"user_id"`
	Value    string          `json:"password"`
	Type     CredentialTypes `json:"type"`
}
