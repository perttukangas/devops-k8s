package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func connectNats() *nats.Conn {
	natsURL := os.Getenv("NATS_URL")
	for {
		nc, err := nats.Connect(natsURL)
		if err == nil {
			fmt.Println("Connected to NATS successfully.")
			return nc
		}
		fmt.Println("Failed to connect to NATS. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

var (
	nc *nats.Conn
)

func main() {
	nc = connectNats()
	defer nc.Close()

	_, err := nc.QueueSubscribe("todo-updates", "todo-group", func(m *nats.Msg) {
		
		webhookURL := os.Getenv("DISCORD_WEBHOOK")
		if webhookURL == "" {
			fmt.Println("DISCORD_WEBHOOK environment variable is not set.")
			return
		}

		message := fmt.Sprintf("Received a todo update: %s", string(m.Data))
		payload := fmt.Sprintf(`{"content": "%s"}`, message)

		if os.Getenv("PRODUCTION") != "true" {
			fmt.Printf("Received a todo update: %s", string(m.Data))
			return
		}

		req, err := http.NewRequest("POST", webhookURL, strings.NewReader(payload))
		if err != nil {
			fmt.Printf("Failed to create request: %v\n", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to send request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Failed to send message to Discord, status code: %d\n", resp.StatusCode)
		} else {
			fmt.Printf("Message sent to Discord successfully: %s\n", message)
		}
	})

	if err != nil {
		fmt.Printf("Failed to subscribe to todo-updates: %v\n", err)
		return
	}

	select {}
}