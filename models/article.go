package models

import "gorm.io/gorm"

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title         string `json:"title" binding:"required" validate:"max=100"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover"`
	CreatedBy     string `json:"created_by" binding:"required"`
	ModifiedBy    string `json:"modified_by"`
	State         int    `json:"state" validate:"oneof'0 1"`
}

func ExistArticleByName(name string) bool {
	var article Article
	db.Select("id").Where("name = ?", name).First(&article)

	return article.ID > 0
}

func ExistArticleByID(id int) (bool, error) {
	var article Article
	err := db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, nil
	}

	if article.ID > 0 {
		return true, nil
	}

	return false, nil
}

func AddArticle(data map[string]interface{}) error {
	article := Article{
		TagID:         data["tag_id"].(int),
		Title:         data["title"].(string),
		Desc:          data["desc"].(string),
		Content:       data["content"].(string),
		CoverImageUrl: data["cover"].(string),
		CreatedBy:     data["created_by"].(string),
		State:         data["state"].(int),
	}

	if err := db.Create(&article).Error; err != nil {
		return err
	}

	return nil
}

func DeleteArticle(id int) error {
	if err := db.Where("id = ?", id).Delete(&Article{}).Error; err != nil {
		return err
	}

	return nil
}

// EditArticle modify a single article
func EditArticle(id int, data interface{}) error {
	if err := db.Model(&Article{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// GetArticle Get a single article based on ID
func GetArticle(id int) (*Article, error) {
	var article Article
	err := db.Where("id = ? AND delete_on = ?", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	err = db.Model(&article).Joins("Tag").Find(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &article, nil
}

func GetArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles []*Article
	err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return articles, nil
}

func GetArticleTotal(maps interface{}) (int64, error) {
	var count int64

	if err := db.Model(&Article{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
