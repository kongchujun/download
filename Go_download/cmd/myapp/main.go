package main

import (
	api "godownload/api"
	config "godownload/configs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	err := config.LoadConfig("/config/config.yaml")
	// err := config.LoadConfig("/Users/kcj/projects/godownload/dev/dcs-download/config/config.yaml")
	if err != nil {
		panic(err)
	}

	router := api.SetupRouter()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// 启动应用程序
	router.Run(":8090")
}
