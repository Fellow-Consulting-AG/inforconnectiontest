package main

import (
	"log"
	"os"
)

var debugMode bool = false
var checkM3 bool = false

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatal("Usage: Infor-test.exe <ionapi-file-path> [--debug] [--check_m3]")
	}

	// Parse arguments for --debug and --check_m3 flags
	ionAPIFile := args[0]
	for _, arg := range args {
		if arg == "--debug" {
			debugMode = true
			log.Println("Debug mode enabled")
		}
		if arg == "--check_m3" {
			checkM3 = true
			log.Println("--check_m3 flag provided")
		}
	}

	log.Printf("Loading ionapi file: %s\n", ionAPIFile)

	// Load the ionAPI file
	ionAPI, err := loadIonAPI(ionAPIFile)
	if err != nil {
		log.Fatalf("Failed to load ionapi file: %v", err)
	}
	log.Println("Successfully loaded ionapi file")

	// Print loaded data if debug mode is enabled
	debugPrint("ION API Gateway URL: %s", ionAPI.IonBaseURL)
	debugPrint("Authorization Server URL: %s", ionAPI.TokenBaseURL)
	debugPrint("Client ID: %s", ionAPI.ClientID)
	debugPrint("Client Secret: %s", ionAPI.ClientSecret)
	debugPrint("Username (SAAK): %s", ionAPI.Username)
	debugPrint("Password (SASK): %s", ionAPI.Password)

	// Check connectivity to ION API Gateway
	if checkConnectivity(ionAPI.IonBaseURL, "ION API Gateway") {
		log.Printf("✅ Successfully connected to ION API Gateway (%s)\n", ionAPI.IonBaseURL)
	} else {
		log.Fatalf("❌ Failed to connect to ION API Gateway (%s)\n", ionAPI.IonBaseURL)
	}

	// Check connectivity to Authorization Server (use the full token URL)
	tokenURL := ionAPI.GetTokenURL() // Use the full URL for checking connectivity
	if checkConnectivity(tokenURL, "Authorization Server") {
		log.Printf("✅ Successfully connected to Authorization Server (%s)\n", tokenURL)
	} else {
		log.Fatalf("❌ Failed to connect to Authorization Server (%s)\n", tokenURL)
	}

	// These connection messages should only be printed once
	log.Println("✅ Connection possible to both ION API Gateway and Authorization Server 💪")

	// Obtain the access token
	token, err := getAccessToken(ionAPI)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}
	log.Println("✅ Connection successful! Access token obtained successfully.")

	// Print the access token if debug mode is enabled
	debugPrint("Access Token: %s", token)

	// If the --check_m3 flag is present, make the M3 API call
	if checkM3 {
		log.Println("Calling M3 API...")
		err := checkM3API(token, ionAPI.IonBaseURL, ionAPI.TenantID)
		if err != nil {
			log.Fatalf("Failed to check M3 API: %v", err)
		}
		log.Println("Successfully called M3 API")
	} else {
		log.Println("No M3 check requested, skipping")
	}

	log.Println("Program finished successfully")
}
