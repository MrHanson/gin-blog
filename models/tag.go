package models

import "gorm.io/gorm"

type Tag struct {
	Model

	Name       string `json:"name" binding:"required" validate:"max=100"`
	CreatedBy  string `json:"created_by" binding:"required"`
	ModifiedBy string `json:"modified_by" validate:"max=100"`
	State      int    `json:"state" validate:"oneof'0 1"`
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ? AND delete_on = ?", name, 0).First(&tag).Error
	if err != nil {
		return false, err
	}

	return tag.ID > 0, nil
}

func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ?AND delete_on = ?", id, 0).First(&tag).Error
	if err != nil {
		return false, err
	}

	return tag.ID > 0, nil
}

func AddTag(name string, state int, createdBy string) error {
	err := db.Create(&Tag{
		Name:       name,
		CreatedBy:  createdBy,
		ModifiedBy: createdBy,
		State:      state,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteTag(id int) error {
	if err := db.Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		return err
	}

	return nil
}

func EditTag(id int, data interface{}) error {
	if err := db.Model(&Tag{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var (
		tags []Tag
		err  error
	)

	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Find(&tags).Offset(pageNum).Limit(pageSize).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

func GetTagTotal(maps interface{}) (int64, error) {
	var count int64
	if err := db.Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
