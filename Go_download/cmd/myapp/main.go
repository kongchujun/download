package main

import (
	"fmt"
	api "godownload/api"
	config "godownload/configs"
	pool "godownload/pool"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	// if you want to change another, you just create struct and implement SFTPConfig
	//err := config.LoadConfig("/config/config.yaml")
	err := config.NewYamlConfig("/Users/kcj/projects/godownload/download/Go_download/configs/config.yaml")
	if err != nil {
		panic(err)
	}
	pool.LoadPool()
	router := api.SetupRouter()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	go func() {
		if err := router.Run(":8090"); err != nil {
			fmt.Println("start server failed: ", err)
		}
	}()
	//monitor signal, so that it can terminal connecting pool
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	// shut down all connection
	pool.SFTPPool.Shutdown()

	time.Sleep(1 * time.Second)
	log.Println("program has been shut down")
}
