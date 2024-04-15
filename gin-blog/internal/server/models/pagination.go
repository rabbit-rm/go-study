package models

import (
	"gorm.io/gorm"
)

func paginationScope(db *gorm.DB) *gorm.DB {
	return db.Offset(1)
}
