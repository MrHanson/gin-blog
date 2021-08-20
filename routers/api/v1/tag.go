package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/unknwon/com"

	"github.com/MrHanson/gin-blog/models"
	"github.com/MrHanson/gin-blog/pkg/app"
	"github.com/MrHanson/gin-blog/pkg/e"
	"github.com/MrHanson/gin-blog/pkg/export"
	"github.com/MrHanson/gin-blog/pkg/logging"
	"github.com/MrHanson/gin-blog/pkg/setting"
	"github.com/MrHanson/gin-blog/pkg/util"
	"github.com/MrHanson/gin-blog/service/tag_service"
)

func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

func AddTag(c *gin.Context) {
	var tag models.Tag
	err := c.ShouldBindJSON(&tag)

	var appG = app.Gin{C: c}

	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      tag.Name,
		CreatedBy: tag.CreatedBy,
		State:     tag.State,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func EditTag(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()

	var tag models.Tag
	tag.ID = id
	err := c.ShouldBindJSON(&tag)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         tag.ID,
		Name:       tag.Name,
		ModifiedBy: tag.ModifiedBy,
		State:      tag.State,
	}

	exist, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()

	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}

func ImportTag(c *gin.Context) {
	appG := app.Gin{C: c}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
