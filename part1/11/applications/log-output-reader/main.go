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
		contentping, errping := os.ReadFile("/usr/src/app/files/pings.txt")
		if err != nil {	
			c.String(500, fmt.Sprintf("Error reading file: %v", err))
			return
		}
		if errping != nil {	
			c.String(500, fmt.Sprintf("Error reading file: %v", errping))
			return
		}
		c.String(200, fmt.Sprintf("%s\n%s", string(content), string(contentping)))
	})

	fmt.Printf("Server started on port %v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}