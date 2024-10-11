package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
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
	randomUUID := uuid.New().String()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		latestResult := fmt.Sprintf("%s: %s", t.Format(time.RFC3339), randomUUID)
		fmt.Println(latestResult)
		err := writeToFile("/usr/src/app/files/timestamp.txt", latestResult)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
		}
	}
}