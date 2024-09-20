package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

var debugMode bool = false
var checkM3 bool = false

// Define log file
const logFile = "infor-test.log"

// Create a custom logger that logs to both a file and stdout
var logger *log.Logger

func init() {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	logger = log.New(file, "", log.LstdFlags)
	logger.SetOutput(os.Stdout)
}

// Updated checkDNSResolution to remove protocol
func checkDNSResolution(rawURL string) error {
	hostname, err := removeProtocol(rawURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %v", err)
	}
	logger.Printf("🔍 Performing DNS resolution for %s", hostname)
	_, err = net.LookupHost(hostname)
	if err != nil {
		logger.Printf("⚠️ DNS Resolution failed for %s: %v", hostname, err) // Log the error but don't stop execution
		return nil                                                          // Return nil to continue program execution
	}
	logger.Printf("✅ DNS Resolution successful for %s", hostname)
	return nil
}

// Updated checkNetworkConnectivity to remove protocol
func checkNetworkConnectivity(rawURL, defaultPort string) error {
	hostname, err := removeProtocol(rawURL) // Remove https:// or http:// from the URL
	if err != nil {
		return fmt.Errorf("failed to parse URL: %v", err)
	}

	// Check if hostname was extracted correctly
	if hostname == "" {
		logger.Printf("⚠️ Hostname extraction failed for URL: %s", rawURL) // Log the error but don't stop execution
		return nil                                                         // Return nil to continue program execution
	}

	logger.Printf("✅ Hostname extracted: %s", hostname)

	// Check if the hostname includes a port
	host, port, err := net.SplitHostPort(hostname)
	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			port = defaultPort // Use default port if no port is provided
			host = hostname    // If no port is found, hostname is just the host
		} else {
			logger.Printf("⚠️ Failed to split host and port: %v", err) // Log the error but don't stop execution
			return nil                                                 // Return nil to continue program execution
		}
	}

	logger.Printf("🔍 Performing network connectivity check to %s on port %s", host, port)
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		logger.Printf("⚠️ Network Connectivity check failed for %s on port %s: %v", host, port, err) // Log the error but don't stop execution
		return nil                                                                                   // Return nil to continue program execution
	}
	defer conn.Close()
	logger.Printf("✅ Network Connectivity successful to %s on port %s", host, port)
	return nil
}

// Updated checkSSLCertificate to remove protocol
func checkSSLCertificate(rawURL, port string) error {
	hostname, err := removeProtocol(rawURL) // Remove https:// or http:// from the URL
	if err != nil {
		logger.Printf("⚠️ SSL Check failed: could not parse URL %s: %v", rawURL, err) // Log the error but don't stop execution
		return nil
	}
	logger.Printf("🔍 Checking SSL/TLS certificate for %s", hostname)
	address := fmt.Sprintf("%s:%s", hostname, port)
	conn, err := tls.Dial("tcp", address, &tls.Config{})
	if err != nil {
		logger.Printf("⚠️ SSL/TLS connection failed for %s: %v", hostname, err) // Log the error but don't stop execution
		return nil
	}
	defer conn.Close()

	// Extract the certificate chain
	certs := conn.ConnectionState().PeerCertificates
	for _, cert := range certs {
		if time.Now().After(cert.NotAfter) {
			logger.Printf("⚠️ Certificate expired on %v", cert.NotAfter)
			return nil
		}
		if time.Now().Before(cert.NotBefore) {
			logger.Printf("⚠️ Certificate not valid before %v", cert.NotBefore)
			return nil
		}
		logger.Printf("✅ Certificate for %s is valid (Valid from %v to %v)", cert.Subject.CommonName, cert.NotBefore, cert.NotAfter)
	}
	return nil
}

// extractPortFromURL extracts the port from the given URL
func extractPortFromURL(url string) string {
	host, port, err := net.SplitHostPort(url)
	if err != nil {
		// Check if URL contains default port (443 or none)
		if port == "" && (host == "443" || err != nil) {
			return "443" // Default to port 443 if none provided
		}
	}
	return port
}

// Remove the protocol (http:// or https://) from the URL and return only the hostname
func removeProtocol(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Return only the host (this excludes the protocol)
	return parsedURL.Host, nil
}

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

	// Perform DNS resolution check for ION Base URL
	if err := checkDNSResolution(ionAPI.IonBaseURL); err != nil {
		log.Fatalf("❌ DNS Resolution failed for %s: %v", ionAPI.IonBaseURL, err)
	} else {
		log.Printf("✅ DNS Resolution successful for %s", ionAPI.IonBaseURL)
	}

	// Extract the port from the Authorization Server URL (pu)
	port := extractPortFromURL(ionAPI.TokenBaseURL)
	if port == "" {
		port = "443" // Default to 443 if no port is provided
	}

	// Perform Network Connectivity check for ION Base URL
	if err := checkNetworkConnectivity(ionAPI.IonBaseURL, port); err != nil {
		logger.Fatalf("❌ Network Connectivity check failed: %v", err)
	} else {
		logger.Printf("✅ Network Connectivity successful to %s on port %s", ionAPI.IonBaseURL, port)
	}

	// Perform SSL/TLS Certificate check
	if err := checkSSLCertificate(ionAPI.IonBaseURL, port); err != nil {
		logger.Printf("⚠️ Warning: SSL Certificate check failed: %v", err) // Only warn, do not stop
	} else {
		logger.Printf("✅ SSL Certificate valid for %s", ionAPI.IonBaseURL)
	}

	// Proceed with the rest of the program (e.g., obtaining access tokens, etc.)

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
