package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	randomUUID := uuid.New().String()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		fmt.Printf("%s: %s\n", t.Format(time.RFC3339), randomUUID)
	}
}