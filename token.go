package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// getAccessToken sends the HTTP request to get the access token
func getAccessToken(api *IonAPI) (string, error) {
	formData := map[string]string{
		"grant_type": "password",
		"username":   api.Username,
		"password":   api.Password,
	}

	form := bytes.NewBufferString("")
	for key, value := range formData {
		form.WriteString(fmt.Sprintf("%s=%s&", key, value))
	}

	// Get the full token URL using the GetTokenURL method
	tokenURL := api.GetTokenURL()
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(form.String()))
	if err != nil {
		return "", err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", api.ClientID, api.ClientSecret)))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get token, status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

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
