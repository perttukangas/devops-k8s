package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	randomUUID := uuid.New().String()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var latestResult string

	go func() {
		for t := range ticker.C {
			latestResult = fmt.Sprintf("%s: %s", t.Format(time.RFC3339), randomUUID)
			fmt.Println(latestResult)
		}
	}()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, latestResult)
	})

	fmt.Printf("Server started on port %v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}