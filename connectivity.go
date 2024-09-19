package main

import (
	"log"
	"net/http"
	"time"
)

// checkConnectivity checks if the application can connect to the provided URL
func checkConnectivity(url string, serviceName string) bool {
	client := &http.Client{
		Timeout: 5 * time.Second, // Set a timeout of 5 seconds
	}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Printf("❌ Failed to create request for %s: %v", serviceName, err)
		return false
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 400 {
		log.Printf("❌ Cannot reach %s. Error: %v, Status Code: %d", serviceName, err, resp.StatusCode)
		return false
	}

	log.Printf("✅ Successfully connected to %s (%s)", serviceName, url)
	return true
}
