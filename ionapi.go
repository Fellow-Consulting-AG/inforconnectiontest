package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// IonAPI structure
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

// Method to construct the full token URL by combining the Base URL and the token path
func (api *IonAPI) GetTokenURL() string {
	return fmt.Sprintf("%s%s", api.TokenBaseURL, api.TokenPath)
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
		return nil, fmt.Errorf("the .ionapi file is missing one or more required fields")
	}

	return &ionAPI, nil
}
