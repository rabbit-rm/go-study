package crud

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type softDeleteUnixTimestamp struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"column:name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt `gorm:"index,softDelete:milli"`
}

func (*softDeleteUnixTimestamp) TableName() string {
	return "t_soft_delete_unix_time"
}

type softDeleteFlag struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"column:name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag"`
}

func (*softDeleteFlag) TableName() string {
	return "t_soft_delete_flag"
}

type softDeleteMixed struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"column:name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}

func (*softDeleteMixed) TableName() string {
	return "t_soft_delete_flag_unix_mixed"
}
