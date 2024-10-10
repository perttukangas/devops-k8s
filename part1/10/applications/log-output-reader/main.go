package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		content, err := os.ReadFile("/usr/src/app/files/timestamp.txt")
		if err != nil {
			c.String(500, fmt.Sprintf("Error reading file: %v", err))
			return
		}
		c.String(200, string(content))
	})

	fmt.Printf("Server started on port %v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}