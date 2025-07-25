package main

import (
	"github.com/gin-gonic/gin"
	"holger-hahn-website/templates"
)

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Home page route
	r.GET("/", func(c *gin.Context) {
		component := templates.Index()
		c.Header("Content-Type", "text/html")
		component.Render(c.Request.Context(), c.Writer)
	})

	r.Run(":8080")
}