package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
	r.GET("/healthz", func(c *gin.Context) {
		_, err := http.Get("http://ping-pong-svc.myapp-namespace.svc.cluster.local/healthz")
		if err != nil {
			c.String(500, fmt.Sprintf("Error fetching health: %v", err))
			return
		}
		c.String(200, "Healthy")
	})

	r.GET("/", func(c *gin.Context) {
		resp, err := http.Get("http://ping-pong-svc.myapp-namespace.svc.cluster.local")
		if err != nil {
			c.String(500, fmt.Sprintf("Error fetching ping: %v", err))
			return
		}
		defer resp.Body.Close()

		ping, err := io.ReadAll(resp.Body)
		if err != nil {
			c.String(500, fmt.Sprintf("Error reading ping response: %v", err))
			return
		}

		content, err := os.ReadFile("/config/information.txt")
		if err != nil {
			c.String(500, fmt.Sprintf("Error reading file: %v", err))
			return
		}

		fileContent := fmt.Sprintf("file content: %v", string(content))
		messageEnv := fmt.Sprintf("env variable: %v", os.Getenv("MESSAGE"))

		c.String(200, fmt.Sprintf("%s\n%s\n%s\n%s", fileContent, messageEnv, string(latestResult), string(ping)))
	})

	fmt.Printf("Server started on port %v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}