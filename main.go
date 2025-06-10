package main

import (
	"net/http"

	"github.com/JxavierP/competitor-chat/components"
	"github.com/gin-gonic/gin"
)

func main()  {
	router := gin.Default()
	router.HTMLRender = &TemplRender{}
	
	router.GET("/", func(c *gin.Context) {
		c.Render(http.StatusOK, &TemplRender{
			 http.StatusOK,
			components.Index(),
		})
	})

	router.Run(":8080") // Start the server on port 8080
}
