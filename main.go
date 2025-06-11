package main

import (
	"net/http"
	"time"

	"github.com/JxavierP/competitor-chat/web/pages"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.HTMLRender = &TemplRender{}

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

	router.Run(":8080") // Start the server on port 8080
}
