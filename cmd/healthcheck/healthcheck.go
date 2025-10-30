package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("Usage: healthcheck <url>")
		return
	}

	url := os.Args[1]

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Health check failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("Health check passed: %s returned status %d\n", url, resp.StatusCode)
		return
	} else {
		log.Printf("Health check failed: %s returned status %d\n", url, resp.StatusCode)
		return
	}
}
