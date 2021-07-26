package api

import (
	"net/http"

	"github.com/MrHanson/gin-blog/models"
	"github.com/MrHanson/gin-blog/pkg/e"
	"github.com/MrHanson/gin-blog/pkg/logging"
	"github.com/MrHanson/gin-blog/pkg/util"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `binding:"required" validate:"max=50"`
	Password string `binding:"required" validate:"max=50"`
}

func authValid(a *auth) bool {
	usernameLen := len(a.Username)
	passwordLen := len(a.Password)

	return usernameLen > 0 && usernameLen <= 50 && passwordLen > 0 && passwordLen <= 50
}

func GetAuth(c *gin.Context) {
	username, _ := c.GetPostForm("username")
	password, _ := c.GetPostForm("password")

	a := auth{Username: username, Password: password}

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS

	isExist := models.CheckAuth(username, password)
	if authValid(&a) && isExist {
		token, err := util.GenerateToken(username, password)
		if err != nil {
			code = e.ERROR_AUTH
			logging.Info(e.GetMsg(code))
		} else {
			data["token"] = token
			code = e.SUCCESS
		}
	} else {
		logging.Warn("账号或密码有误")
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
