package main

import (
	"fmt"
	"net/http"

	"github.com/MrHanson/gin-blog/models"
	"github.com/MrHanson/gin-blog/pkg/logging"
	"github.com/MrHanson/gin-blog/pkg/setting"
	"github.com/MrHanson/gin-blog/routers"
)

func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()

	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
