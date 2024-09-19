package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Global variable to control check_m3 functionality
var checkM3 bool = false

// Function to check the /M3/m3api-rest/application.wadl endpoint
func checkM3API(token string, baseURL string, tenantID string) error {
	// Construct the full URL
	apiEndpoint := fmt.Sprintf("%s/%s/M3/m3api-rest/application.wadl", baseURL, tenantID)

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return err
	}

	// Add Bearer token to the Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("accept", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Read and print the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("M3 API Response: %s\n", body)
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Usage: Infor-test.exe <ionapi-file-path> [--debug] [--check_m3]")
	}

	// Check for --debug flag and --check_m3 flag
	if len(args) > 1 && args[1] == "--debug" {
		debugMode = true
	}
	if len(args) > 2 && args[2] == "--check_m3" {
		checkM3 = true
	}

	ionAPIFile := args[0] // This should correctly point to the ionapi file

	// Load the ionAPI file
	ionAPI, err := loadIonAPI(ionAPIFile)
	if err != nil {
		log.Fatalf("Failed to load ionapi file: %v", err)
	}

	// Obtain the access token
	token, err := getAccessToken(ionAPI)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	// Check for --check_m3 flag and make the M3 API request
	if checkM3 {
		err := checkM3API(token, ionAPI.IonBaseURL, ionAPI.TenantID)
		if err != nil {
			log.Fatalf("Failed to check M3 API: %v", err)
		}
	}
}
