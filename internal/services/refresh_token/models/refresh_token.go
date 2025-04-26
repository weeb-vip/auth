package models

import "github.com/weeb-vip/auth/internal/db"

type RefreshToken struct {
	db.BaseModel
	UserID string `gorm:"column:user_id;type:varchar(36);not null"`
	Token  string `gorm:"column:token;type:varchar(36);not null"`
	Expiry int64  `gorm:"column:expiry;type:int;not null"`
}
