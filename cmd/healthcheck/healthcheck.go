package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: healthcheck <url>")
	}

	url := os.Args[1]

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Health check failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("Health check passed: %s returned status %d\n", url, resp.StatusCode)
		os.Exit(0)
	} else {
		fmt.Printf("Health check failed: %s returned status %d\n", url, resp.StatusCode)
		os.Exit(1)
	}
}
