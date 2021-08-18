package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/unknwon/com"

	"github.com/MrHanson/gin-blog/models"
	"github.com/MrHanson/gin-blog/pkg/e"
	"github.com/MrHanson/gin-blog/pkg/setting"
	"github.com/MrHanson/gin-blog/pkg/util"
)

func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state = -1
	arg := c.Query("state")
	if arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	data["lists"] = models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func AddTag(c *gin.Context) {
	var tag models.Tag
	err := c.ShouldBindJSON(&tag)

	name := tag.Name
	code := e.INVALID_PARAMS
	if err == nil {
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, tag.State, tag.CreatedBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	var tag models.Tag
	err := c.ShouldBindJSON(&tag)
	code := e.INVALID_PARAMS
	if err == nil {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.EditTag(id, tag)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	code := e.INVALID_PARAMS
	if id > 0 {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
