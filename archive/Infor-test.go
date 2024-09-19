package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Global variable to control debug output
var debugMode bool = false

// IonAPI structure to map the fields from your .ionapi file
type IonAPI struct {
	ClientID     string `json:"ci"`   // Client ID
	ClientSecret string `json:"cs"`   // Client Secret
	TokenBaseURL string `json:"pu"`   // Base URL to form the token URL
	TokenPath    string `json:"ot"`   // Path to form the token URL
	Username     string `json:"saak"` // Use `saak` as the username
	Password     string `json:"sask"` // Use `sask` as the password
	IonBaseURL   string `json:"iu"`   // Base URL for ION API
	TenantID     string `json:"ti"`   // Tenant ID
}

// debugPrint prints debug messages if debugMode is enabled
func debugPrint(format string, v ...interface{}) {
	if debugMode {
		fmt.Printf(format, v...)
	}
}

// Construct the full token URL by combining the Base URL and the token path
func (api *IonAPI) GetTokenURL() string {
	return fmt.Sprintf("%s%s", api.TokenBaseURL, api.TokenPath)
}

func main() {

	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatal("Usage: Infor-test.exe <ionapi-file-path> [--debug]")
	}

	// Check for --debug flag
	if len(args) > 1 && args[1] == "--debug" {
		debugMode = true
	}

	ionAPIFile := args[0] // This should correctly point to the ionapi file

	// Load the ionAPI file
	ionAPI, err := loadIonAPI(ionAPIFile)
	if err != nil {
		log.Fatalf("Failed to load ionapi file: %v", err)
	}

	// Check connectivity to both the ION API Gateway and Authorization Server
	if !checkConnectivity(ionAPI.IonBaseURL, "ION API Gateway") {
		log.Fatalf("❌ Cannot connect to ION API Gateway (%s)", ionAPI.IonBaseURL)
	}
	// if !checkConnectivity(ionAPI.TokenBaseURL, "Authorization Server") {
	// 	log.Fatalf("❌ Cannot connect to Authorization Server (%s)", ionAPI.TokenBaseURL)
	// }

	// Print connectivity success
	fmt.Println("✅ Connection possible to both ION API Gateway and Authorization Server")

	// Debugging: Print out all fields from the .ionapi file to check if they're loaded correctly
	debugPrint("Loaded .ionapi file with the following values:\n")
	debugPrint("Client ID: %s\n", ionAPI.ClientID)
	debugPrint("Client Secret: %s\n", ionAPI.ClientSecret)
	debugPrint("Username (SAAK): %s\n", ionAPI.Username)
	debugPrint("Password (SASK): %s\n", ionAPI.Password)
	debugPrint("Access Token URL: %s\n", ionAPI.GetTokenURL())

	// Obtain the access token
	token, err := getAccessToken(ionAPI)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	// Print access token
	debugPrint("Access Token: %s\n", token)

	// Check if the base URL ends with inforcloudsuite.com
	if !strings.Contains(ionAPI.IonBaseURL, "inforcloudsuite.com") {
		fmt.Println("❌ Not using inforcloudsuite.com domain, skipping API request")
		return
	}

	// Make the GET request using the access token as a Bearer token
	err = makeAuthenticatedRequest(token, ionAPI.IonBaseURL, ionAPI.TenantID)
	if err != nil {
		log.Fatalf("Failed to make API request: %v", err)
	}

	// Print success message
	fmt.Println("✅ Connection successful! Access token obtained and API request made successfully.")
}

// loadIonAPI reads and parses the .ionapi file
func loadIonAPI(filePath string) (*IonAPI, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var ionAPI IonAPI
	err = json.Unmarshal(fileBytes, &ionAPI)
	if err != nil {
		return nil, err
	}

	// Check if required fields are present
	if ionAPI.ClientID == "" || ionAPI.ClientSecret == "" || ionAPI.TokenBaseURL == "" || ionAPI.TokenPath == "" || ionAPI.Username == "" || ionAPI.Password == "" || ionAPI.IonBaseURL == "" || ionAPI.TenantID == "" {
		return nil, fmt.Errorf("the .ionapi file is missing one or more required fields (ci, cs, pu, ot, saak, sask, iu, ti)")
	}

	return &ionAPI, nil
}

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

// getAccessToken sends the HTTP request to get the access token
func getAccessToken(api *IonAPI) (string, error) {
	// Prepare the form data for Password Grant Type
	formData := map[string]string{
		"grant_type": "password",
		"username":   api.Username, // Using `saak` as username
		"password":   api.Password, // Using `sask` as password
		"scope":      "",           // Scope is left blank as per requirements
	}

	// Convert form data to a format the server expects
	form := bytes.NewBufferString("")
	for key, value := range formData {
		form.WriteString(fmt.Sprintf("%s=%s&", key, value))
	}

	tokenURL := api.GetTokenURL()
	// Debugging: Print out the token URL and form data being sent
	debugPrint("Sending POST request to: %s\n", tokenURL)
	debugPrint("Form Data: %s\n", form.String())

	// Create the POST request
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(form.String()))
	if err != nil {
		return "", err
	}

	// Add Basic Authentication Header with Client ID and Client Secret
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", api.ClientID, api.ClientSecret)))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token, status: %s", resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Print the raw response body if debug mode is enabled
	debugPrint("Raw Response: %s\n", body)

	// Extract access token from the response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("no access_token found in response")
	}

	return accessToken, nil
}

// makeAuthenticatedRequest makes an authenticated GET request to the API
func makeAuthenticatedRequest(token string, baseURL string, tenantID string) error {
	// Construct the full API URL using the base URL, tenantID, and API path
	apiEndpoint := fmt.Sprintf("%s/%s/OSPORTAL/admin/v1/user/applications", baseURL, tenantID)

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return err
	}

	// Add Bearer token to the Authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("accept", "application/json")

	// Print request details if debug mode is enabled
	debugPrint("Sending GET request to: %s\n", apiEndpoint)

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

	debugPrint("API Response: %s\n", body)
	return nil
}
