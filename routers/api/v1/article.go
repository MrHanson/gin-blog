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

func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	code := e.INVALID_PARAMS
	if id > 0 {
		if models.ExistArticleByID(id) {
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	}

	data := make(map[string]interface{})

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func GetArticles(c *gin.Context) {
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	title := c.Query("title")
	if title != "" {
		maps["title"] = title
	}

	var state = -1
	arg := c.Query("state")
	if arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"] = models.GetArticleTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

func AddArticle(c *gin.Context) {
	var article models.Article
	err := c.ShouldBindJSON(&article)

	tagId := article.TagID
	title := article.Title
	desc := article.Desc
	content := article.Content
	createdBy := article.CreatedBy
	state := article.State

	code := e.INVALID_PARAMS
	if err == nil {
		if models.ExistTagByID(tagId) {
			data := make(map[string]interface{})
			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
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

func EditArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	var article models.Article
	err := c.ShouldBindJSON(&article)
	code := e.INVALID_PARAMS
	if err == nil {
		if models.ExistArticleByID(id) {
			if models.ExistTagByID(article.TagID) {
				data := make(map[string]interface{})
				data["tag_id"] = article.TagID
				data["title"] = article.Title
				data["desc"] = article.Desc
				data["content"] = article.Content
				data["modified_by"] = article.ModifiedBy

				models.EditArticle(id, data)
			} else {
				code = e.ERROR_NOT_EXIST_TAG
			}
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	code := e.INVALID_PARAMS
	if id > 0 {
		code = e.SUCCESS
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
