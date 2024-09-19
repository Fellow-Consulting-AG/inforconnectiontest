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
)

// IonAPI structure to map the fields from your .ionapi file
type IonAPI struct {
	ClientID     string `json:"ci"`   // Client ID
	ClientSecret string `json:"cs"`   // Client Secret
	TokenBaseURL string `json:"pu"`   // Base URL to form the token URL
	TokenPath    string `json:"ot"`   // Path to form the token URL
	Username     string `json:"saak"` // Use `saak` as the username
	Password     string `json:"sask"` // Use `sask` as the password
}

// Construct the full token URL by combining the Base URL and the token path
func (api *IonAPI) GetTokenURL() string {
	return fmt.Sprintf("%s%s", api.TokenBaseURL, api.TokenPath)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: Infor-test.exe <ionapi-file-path>")
	}

	ionAPIFile := os.Args[1]
	ionAPI, err := loadIonAPI(ionAPIFile)
	if err != nil {
		log.Fatalf("Failed to load ionapi file: %v", err)
	}

	// Debugging: Print out all fields from the .ionapi file to check if they're loaded correctly
	fmt.Println("Loaded .ionapi file with the following values:")
	fmt.Printf("Client ID: %s\n", ionAPI.ClientID)
	fmt.Printf("Client Secret: %s\n", ionAPI.ClientSecret)
	fmt.Printf("Username (SAAK): %s\n", ionAPI.Username)
	fmt.Printf("Password (SASK): %s\n", ionAPI.Password)
	fmt.Printf("Access Token URL: %s\n", ionAPI.GetTokenURL())

	// Obtain the access token
	token, err := getAccessToken(ionAPI)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	// Print access token
	fmt.Printf("Access Token: %s\n", token)

	// Print success message
	fmt.Println("✅ Connection successful! Access token obtained successfully.")
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
	if ionAPI.ClientID == "" || ionAPI.ClientSecret == "" || ionAPI.TokenBaseURL == "" || ionAPI.TokenPath == "" || ionAPI.Username == "" || ionAPI.Password == "" {
		return nil, fmt.Errorf("the .ionapi file is missing one or more required fields (ci, cs, pu, ot, saak, sask)")
	}

	return &ionAPI, nil
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
	fmt.Printf("Sending POST request to: %s\n", tokenURL)
	fmt.Printf("Form Data: %s\n", form.String())

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

	// Debugging: Print the raw response body
	fmt.Printf("Raw Response: %s\n", body)

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
