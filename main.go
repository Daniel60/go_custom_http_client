package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Daniel60/go_custom_http_client/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
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

	// Endpoint /todos com chamada para API externa
	g.GET("/todos", func(ctx *gin.Context) {
		// Gerar um valor aleatório para o header
		randomHeader := fmt.Sprintf("random-%d", rand.Intn(1000))

		// Criar uma requisição para a API externa
		externalAPIURL := "https://jsonplaceholder.typicode.com/todos/1"
		req, err := http.NewRequest("GET", externalAPIURL, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create request",
			})
			return
		}

		// Adicionar headers à requisição
		req.Header.Add("X-Random-Header", randomHeader)
		req.Header.Add("X-Intermediario", "intermediario")
		req.Header.Add("Authorization", "Bearer 1234567890234567890")

		ctx.Set("external_request_headers", req.Header)

		// Fazer a requisição para a API externa
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to call external API",
			})
			return
		}
		defer resp.Body.Close()

		// Copiar os cabeçalhos da resposta externa para a resposta do cliente
		for key, values := range resp.Header {
			for _, value := range values {
				ctx.Writer.Header().Add(key, value)
			}
		}

		// Definir o status code da resposta externa
		ctx.Status(resp.StatusCode)

		// Copiar o corpo da resposta externa diretamente para o cliente
		_, err = io.Copy(ctx.Writer, resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to copy external API response",
			})
			return
		}
	})

	g.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    "PATH_NOT_FOUND",
			"message": "Path not Found",
		})
	})
	g.Run(":8080")
}
