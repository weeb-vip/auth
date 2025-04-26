package models

import (
	"github.com/weeb-vip/auth/internal/db"
)

type Session struct {
	db.BaseModel
	UserID    string `column:"user_id"`
	IPAddress string `column:"ip_address"`
	UserAgent string `column:"user_agent"`
	Token     string `column:"token"`
}
