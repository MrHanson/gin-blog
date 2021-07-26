package jwt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/MrHanson/gin-blog/pkg/e"
	"github.com/MrHanson/gin-blog/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.SUCCESS
		token := c.GetHeader("token")
		if token == "" {
			code = e.ERROR_AUTH
		} else {
			clamins, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > clamins.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": make(map[string]string),
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
