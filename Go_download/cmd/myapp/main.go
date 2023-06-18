package main

import (
	api "godownload/api"
	config "godownload/configs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// if you want to change another, you just create struct and implement SFTPConfig
	//err := config.LoadConfig("/config/config.yaml")
	err := config.NewYamlConfig("/Users/kcj/projects/godownload/download/Go_download/configs/config.yaml")
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
