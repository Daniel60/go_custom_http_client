package main

import (
	"net/http"

	"github.com/Daniel60/go_custom_http_client/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(logger.CustomGinLogger())

	g.GET("/v1/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "UP",
		})
	})
	g.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "UP",
		})
	})
	g.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "UP",
		})
	})

	g.GET("/films", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "sups",
		})
	})

	g.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    "PATH_NOT_FOUND",
			"message": "Path not Found",
		})
	})
	g.Run(":8080")
}
