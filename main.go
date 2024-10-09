package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

		// Ler a resposta da API externa
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read external API response",
			})
			return
		}

		// Converter a resposta para um mapa
		var responseBody map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to parse external API response",
			})
			return
		}

		// Retornar a resposta da API externa junto com os headers adicionados
		ctx.JSON(http.StatusOK, gin.H{
			"message":       "sups",
			"external_api":  responseBody,
			"random_header": randomHeader,
			"intermediario": "intermediario",
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
