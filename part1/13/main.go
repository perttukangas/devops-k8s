package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func fetchImage() {
	filePath := "/usr/src/app/ui/images/image.jpg"
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		if time.Since(fileInfo.ModTime()) < 60*time.Minute {
			fmt.Println("Image is less than 60 minutes old, skipping fetch.")
			return
		}
	} else if !os.IsNotExist(err) {
		fmt.Println("Error checking file:", err)
		return
	}

	resp, err := http.Get("https://picsum.photos/1200")
	if err != nil {
		fmt.Println("Error fetching image:", err)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}
	fmt.Println("Image fetched and saved successfully.")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	go func() {
		for {
			fetchImage()
			fileInfo, err := os.Stat("/usr/src/app/ui/images/image.jpg")
			if err == nil {
				timeSinceMod := time.Since(fileInfo.ModTime())
				if timeSinceMod < 60*time.Minute {
					sleepDuration := 60*time.Minute - timeSinceMod
					fmt.Printf("Sleeping for %v\n", sleepDuration)
					time.Sleep(sleepDuration)
				}
			} else {
				fmt.Println("Sleeping for 60 minutes")
				time.Sleep(60 * time.Minute)
			}
		}
	}()

	r := gin.Default()

	r.Static("/", "./ui")

	fmt.Printf("Server started on port %v\n", port)
	r.Run(fmt.Sprintf(":%v", port))
}
