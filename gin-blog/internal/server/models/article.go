package models

import (
	"blog/internal/server/db"

	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	TagID       uint   `gorm:"column:tag_id;comment:标签ID"`
	Tag         Tag    `gorm:"foreignKey:TagID"`
	Title       string `gorm:"column:title;comment:文章标题"`
	Description string `gorm:"column:description;comment:简述"`
	Content     string `gorm:"column:content;comment:内容"`
	CreateBy    string `gorm:"column:create_by;comment:创建人"`
	ModifiedBy  string `gorm:"column:modified_by;comment:修改人"`
	State       uint   `gorm:"column:state;default:1;comment:状态 0 禁用、1 启用"`
}

func ExistArticleByID(id uint64) bool {
	var article Article
	db.MySQL().Select("id").Where("id = ?", id).First(&article)
	return article.ID > 0
}

func GetArticleTotal(condition map[string]string) (count int64) {
	db.MySQL().Model(&Article{}).Where(condition).Count(&count)
	return
}

func GetArticleList(condition map[string]string, pageSize int, pageIndex int) ([]Article, int64) {
	var articles []Article
	offset := (pageIndex - 1) * pageSize
	db.MySQL().Preload("Tag").Where(condition).Offset(offset).Limit(pageSize).Find(&articles)
	return articles, int64(len(articles))
}

func GetArticle(id uint64) Article {
	var article Article
	db.MySQL().Preload("Tag").Where("`id` = ?", id).First(&article)
	return article
}

func EditArticle(id uint64, data Article) bool {
	tx := db.MySQL().Model(&Article{}).Where("id = ?", id).Updates(&data)
	return tx.RowsAffected == 1
}

func AddArticle(article Article) bool {
	tx := db.MySQL().Create(&article)
	return tx.RowsAffected == 1
}

func DeleteArticle(id uint64) bool {
	tx := db.MySQL().Where("id = ?", id).Delete(&Article{})
	return tx.RowsAffected == 1
}
