package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// checkM3API calls the /M3/m3api-rest/application.wadl endpoint
func checkM3API(token string, baseURL string, tenantID string) error {
	apiEndpoint := fmt.Sprintf("%s/%s/M3/m3api-rest/v2/execute/CMS535MI/FpwVersion?dateformat=YMD8&excludeempty=false&righttrim=true&format=PRETTY&extendedresult=false", baseURL, tenantID)

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("M3 API Response: %s\n", body)
	return nil
}
