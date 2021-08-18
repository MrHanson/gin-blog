package routers

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/MrHanson/gin-blog/middleware/jwt"
	"github.com/MrHanson/gin-blog/pkg/setting"
	"github.com/MrHanson/gin-blog/pkg/upload"
	"github.com/MrHanson/gin-blog/routers/api"
	v1 "github.com/MrHanson/gin-blog/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	corsMw := cors.New(cors.Config{
		//准许跨域请求网站,多个使用,分开,限制使用*
		AllowOrigins: []string{"*"},
		//准许使用的请求方式
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		//准许使用的请求表头
		AllowHeaders: []string{"Origin", "token", "Content-Type"},
		//显示的请求表头
		ExposeHeaders: []string{"Content-Type"},
		//凭证共享,确定共享
		AllowCredentials: true,
		//容许跨域的原点网站,可以直接return true就万事大吉了
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		//超时时间设定
		MaxAge: 5 * time.Hour,
	})
	r.Use(corsMw)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	r.POST("/auth", api.GetAuth)
	r.POST("/upload", api.UploadImage)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		apiv1.GET("/article/:id", v1.GetArticle)
		apiv1.GET("/articles", v1.GetArticles)
		apiv1.POST("/article", v1.AddArticle)
		apiv1.PUT("/article/:id", v1.EditArticle)
		apiv1.DELETE("/article/:id", v1.DeleteArticle)
	}

	return r
}
