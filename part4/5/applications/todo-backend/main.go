package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func connectDb() (*sql.DB) {
	dbName := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	for {
		db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		  host, dbPort, user, password, dbName))
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			fmt.Println("Connected to database successfully.")
			return db
		}
		fmt.Println("Failed to connect to database. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

func initDb(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		todo TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT FALSE
	);`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Printf("Error creating table: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Table 'todos' created or already exists.")
}

type Todo struct {
	ID   int    `json:"id"`
	Todo string `json:"todo"`
	Done bool   `json:"done"`
}

func getTodos(c *gin.Context) {
	rows, err := db.Query("SELECT id, todo, done FROM todos ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var todos = []Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Todo, &todo.Done); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context) {
	var newTodo Todo
	if err := c.BindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(newTodo.Todo) > 140 {
		fmt.Println("Todo length validation failed (length):", newTodo.Todo)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todo cannot be longer than 140 characters"})
		return
	}

	var id int
	err := db.QueryRow("INSERT INTO todos (todo) VALUES ($1) RETURNING id", newTodo.Todo).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newTodo.ID = id
	newTodo.Done = false

	fmt.Println("Successfully added new todo:", newTodo)
	c.JSON(http.StatusCreated, newTodo)
}

func putTodo(c *gin.Context) {
	id := c.Param("id")
	var todo Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(todo.Todo) > 140 {
		fmt.Println("Todo length validation failed (length):", todo.Todo)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todo cannot be longer than 140 characters"})
		return
	}

	err := db.QueryRow("UPDATE todos SET todo = $1, done = $2 WHERE id = $3", todo.Todo, todo.Done, id).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Successfully updated todo:", todo)
	c.JSON(http.StatusOK, todo)
}

var (
	db *sql.DB
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	db = connectDb()
	initDb(db)

	r := gin.Default()

	fmt.Printf("Init healthz...")
	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "Healthy")
	})

	fmt.Printf("Init get...")
	r.GET("/", getTodos)

	fmt.Printf("Init post...")
	r.POST("/", createTodo)

	fmt.Printf("Init put...")
	r.PUT("/:id", putTodo)

	fmt.Printf("Server started on port %v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}
