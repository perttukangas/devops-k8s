package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func writeToFile(filename, data string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
			return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	return err
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	r := gin.Default()

	counter := 0
	r.GET("/", func(c *gin.Context) {
		counter++
		err := writeToFile("/usr/src/app/files/pings.txt", fmt.Sprintf("%d", counter))
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
		}
		c.String(200, fmt.Sprintf("pong %d", counter))
	})

	fmt.Printf("Server started on port %v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}