package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new gin router
	r := gin.Default()

	// Define a simple route for test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Start the server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
