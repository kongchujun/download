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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	//add promethues monitor
	// p := ginprometheus.NewPrometheus("my_download")
	// p.Use(router)

	requestCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Total number of HTTP requests",
	})
	requestDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "HTTP request duration in seconds",
	})
	requestSize := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "http_request_size_bytes",
		Help: "HTTP request size in bytes",
	})
	responseSize := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "http_response_size_bytes",
		Help: "HTTP response size in bytes",
	})
	// register all indictor
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(requestSize)
	prometheus.MustRegister(responseSize)
	// middleware to monitor
	router.Use(func(ctx *gin.Context) {
		start := time.Now()
		requestCounter.Inc()
		requestSize.Observe(float64(ctx.Request.ContentLength))
		ctx.Next()

		responseSize.Observe(float64(ctx.Writer.Size()))
		duration := time.Since(start).Seconds()
		requestDuration.Observe(duration)
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
