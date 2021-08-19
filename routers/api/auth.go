package api

import (
	"net/http"

	"github.com/MrHanson/gin-blog/pkg/app"
	"github.com/MrHanson/gin-blog/pkg/e"
	"github.com/MrHanson/gin-blog/pkg/util"
	"github.com/MrHanson/gin-blog/service/auth_service"
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
	appG := app.Gin{C: c}
	username, _ := c.GetPostForm("username")
	password, _ := c.GetPostForm("password")

	a := auth{Username: username, Password: password}
	if !authValid(&a) {
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
