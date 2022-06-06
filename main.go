package main

import (
	"awesomeProject/handler"
	"awesomeProject/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowMethods("GET", "POST")
	r.Use(cors.New(corsConfig))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	healthCheckService := service.NewHealthCheckService()
	healthCheckHandler := handler.NewUploadHandler(healthCheckService)
	r.POST("/upload", healthCheckHandler.Upload)
	r.Run()
}
