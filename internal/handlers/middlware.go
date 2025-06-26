package handlers

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/genai"
)

func GeminiMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// This middleware can be used to handle Gemini-specific logic
		// For now, it just passes the request to the next handler
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  "",
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to create Gemini client"})
			return
		}

		if client.ClientConfig().APIKey == "" {
			ctx.JSON(400, gin.H{"error": "API key is required for Gemini client"})
			return
		}

		ctx.Next()
	}
}
