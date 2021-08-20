package tag_service

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/MrHanson/gin-blog/models"
	"github.com/MrHanson/gin-blog/pkg/export"
	"github.com/MrHanson/gin-blog/pkg/gredis"
	"github.com/MrHanson/gin-blog/pkg/logging"
	"github.com/MrHanson/gin-blog/service/cache_service"
	"github.com/xuri/excelize/v2"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}

	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int64, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	cache := cache_service.Tag{
		State: t.State,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	gredis.Set(key, tags, 3600)
	return tags, nil
}

const EXPORT_SHEET_NAME = "标签信息"

func getRowAxis(index int64) string {
	ret := fmt.Sprintf("%c", 65+index)
	if index > 25 {
		return ret + getRowAxis(index-25)
	}

	return ret
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	index := f.NewSheet(EXPORT_SHEET_NAME)

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	// set title
	f.SetSheetRow(EXPORT_SHEET_NAME, "A1", titles)

	// set data content
	for j, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}

		f.SetSheetRow(EXPORT_SHEET_NAME, getRowAxis(int64(j)), values)

	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	time := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-" + time + ".xlsx"
	fullPath := export.GetExcelFullPath() + filename
	if err := f.SaveAs(fullPath); err != nil {
		return "", err
	}

	return filename, nil
}

func (t *Tag) Import(r io.Reader) error {
	xlsx, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows, err := xlsx.GetRows(EXPORT_SHEET_NAME)
	if err != nil {
		return err
	}
	for irow, row := range rows {
		if irow < 0 {
			continue
		}
		var data []string
		data = append(data, row...)

		models.AddTag(data[1], 1, data[2])
	}

	return nil
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}
