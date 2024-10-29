package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	r := gin.Default()

	dbName := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")

	var db *sql.DB
	var err error
	for {
		db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, dbPort, user, password, dbName))
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		fmt.Println("Failed to connect to database. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	fmt.Printf("Connected to database %s successfully\n", dbName)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS counter (id SERIAL PRIMARY KEY, count INT NOT NULL)`)
	if err != nil {
		fmt.Printf("Failed to create counter table: %v\n", err)
		return
	}

	_, err = db.Exec(`INSERT INTO counter (count) VALUES (0) ON CONFLICT DO NOTHING`)
	if err != nil {
		fmt.Printf("Failed to initialize counter: %v\n", err)
		return
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "Healthy")
	})	

	r.GET("/", func(c *gin.Context) {
		_, err := db.Exec(`UPDATE counter SET count = count + 1 WHERE id = 1`)
		if err != nil {
			fmt.Printf("Failed to update counter: %v\n", err)
			c.String(500, "Internal Server Error")
			return
		}

		var counter int
		err = db.QueryRow(`SELECT count FROM counter WHERE id = 1`).Scan(&counter)
		if err != nil {
			fmt.Printf("Failed to get counter: %v\n", err)
			c.String(500, "Internal Server Error")
			return
		}

		c.String(200, fmt.Sprintf("pong %d", counter))
	})

	fmt.Printf("Server started on port %v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}