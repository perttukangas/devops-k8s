package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	r := gin.Default()

	r.Static("/", "./ui")

	fmt.Printf("Server started on port %v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}
