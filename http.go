package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Decodes the JSON response from the HTTP request into the provided struct
func decodeResponse(resp *http.Response, v any) error {
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed to decode response: " + resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

// Sends an HTTP GET request to the specified URL with the Authorization header set
func sendRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed: " + resp.Status)
	}

	return resp, nil
}
