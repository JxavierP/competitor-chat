package main

import (
	"context"
	"net/http"
	"time"

	"github.com/JxavierP/competitor-chat/internal/handlers"
	"github.com/JxavierP/competitor-chat/web/pages"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

func main() {
	router := gin.Default()
	router.HTMLRender = &TemplRender{}
	gemini, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: "AIzaSyDIggMrw4BOpML-6KOyrsl6ONVotE397WM",
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		println("Error creating Gemini client:", err)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "https://hoppscotch.io"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		c.Render(http.StatusOK, &TemplRender{
			http.StatusOK,
			pages.Index(),
		})
	})
	router.POST("/chat", handlers.PromptHandler)
	router.GET("/chat", handlers.ResponseHandler(gemini))

	router.Run(":8080") // Start the server on port 8080
}
