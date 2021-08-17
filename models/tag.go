package models

type Tag struct {
	Model

	Name       string `json:"name" binding:"required" validate:"max=100"`
	CreatedBy  string `json:"created_by" binding:"required"`
	ModifiedBy string `json:"modified_by" validate:"max=100"`
	State      int    `json:"state" validate:"oneof'0 1"`
}

func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name = ?", name).First(&tag)

	return tag.ID > 0
}

func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)

	return tag.ID > 0
}

func AddTag(name string, state int, createdBy string) bool {
	db.Create(&Tag{
		Name:       name,
		CreatedBy:  createdBy,
		ModifiedBy: createdBy,
		State:      state,
	})

	return true
}

func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})

	return true
}

func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id = ?", id).Updates(data)

	return true
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tags []Tag) {
	db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)

	return
}

func GetTagTotal(maps interface{}) (count int) {
	c := int64(count)
	db.Model(&Tag{}).Where(maps).Count(&c)

	return
}
