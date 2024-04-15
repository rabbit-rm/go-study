package models

import (
	"blog/internal/server/db"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name       string `gorm:"column:name;comment:标签名称"`
	CreateBy   string `gorm:"column:create_by;comment:创建人"`
	ModifiedBy string `gorm:"column:modified_by;comment:修改人"`
	State      uint   `gorm:"column:state;default:1;comment:状态 0 禁用、1 启用"`
}

func ExistTagById(id uint64) bool {
	var tag Tag
	db.MySQL().Model(&Tag{}).Select("id").Where("id = ?", id).First(&tag)
	return tag.ID > 0
}

func GetTag(id uint64) (*Tag, error) {
	var tag Tag
	err := db.MySQL().Model(&Tag{}).Where("id = ?", id).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}
func AddTag(tag *Tag) error {
	return db.MySQL().Create(tag).Error
}
func EditTag(id uint64, tag *Tag) error {
	return db.MySQL().Model(&Tag{}).Where("id = ?", id).Updates(tag).Error
}

func DeleteTag(id uint64) error {
	return db.MySQL().Model(&Tag{}).Where("id = ?", id).Delete(&Tag{}).Error
}
