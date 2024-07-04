package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/v1/health"},
	}))

	g.GET("/v1/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "UP",
		})
	})

	g.Run(":8080")
}
