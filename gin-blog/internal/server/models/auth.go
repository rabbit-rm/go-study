package models

import (
	"blog/internal/server/db"

	"gorm.io/gorm"
)

type Auth struct {
	gorm.Model
	Username string `gorm:"COLUMN:username;COMMENT:账号"`
	Password string `gorm:"COLUMN:password;COMMENT:密码"`
}

func GetAuthByUsername(username string) (*Auth, error) {
	var auth Auth
	tx := db.MySQL().Model(&Auth{}).Where("username = ?", username).Find(&auth)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &auth, nil
}
