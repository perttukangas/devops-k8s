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
		todo TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Printf("Error creating table: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Table 'todos' created or already exists.")
}

type Todo struct {
	Todo string `json:"todo"`
}

func getTodos(c *gin.Context) {
	rows, err := db.Query("SELECT todo FROM todos ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var todos = []Todo{}
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.Todo); err != nil {
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

	_, err := db.Exec("INSERT INTO todos (todo) VALUES ($1)", newTodo.Todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTodo)
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

	fmt.Printf("Init get...")
	r.GET("/", getTodos)

	fmt.Printf("Init post...")
	r.POST("/", createTodo)

	fmt.Printf("Server started on port %v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}
