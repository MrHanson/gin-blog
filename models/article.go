package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title      string `json:"title" binding:"required" validate:"max=100"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by" binding:"required"`
	ModifiedBy string `json:"modified_by" validate:"max=100"`
	State      int    `json:"state" validate:"oneof'0 1"`
}

func (article *Article) BeforeCreate(tx *gorm.DB) (err error) {
	article.CreatedOn = int(time.Now().Unix())

	return nil
}

func (article *Article) BeforeUpdate(tx *gorm.DB) (err error) {
	article.ModifiedOn = int(time.Now().Unix())

	return nil
}

func ExistArticleByName(name string) bool {
	var article Article
	db.Select("id").Where("name = ?", name).First(&article)

	return article.ID > 0
}

func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(&article)

	return article.ID > 0
}

func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	})

	return true
}

func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(&Article{})

	return true
}

func EditArticle(id int, data interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)

	return true
}

func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).First(&article)
	db.Model(&article).Joins("Tag").Find(&article)

	return
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

func GetArticleTotal(maps interface{}) (count int) {
	c := int64(count)
	db.Model(&Article{}).Where(maps).Count(&c)

	return
}