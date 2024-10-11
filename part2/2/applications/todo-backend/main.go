package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Todo string `json:"todo"`
}

func getTodos(c *gin.Context) {
	c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.BindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todos = append(todos, newTodo)
	c.JSON(http.StatusCreated, newTodo)
}

var (
	todos = []Todo{}
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	r := gin.Default()

	fmt.Printf("Init get...")
	r.GET("/", getTodos)

	fmt.Printf("Init post...")
	r.POST("/", createTodo)

	fmt.Printf("Server started on port %v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}
