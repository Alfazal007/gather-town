package usertests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestSignupEndpoint(t *testing.T) {
	// Prepare the request payload
	payload := map[string]string{
		"username": "someone1",
		"email":    "someone1@somewhere.com",
		"password": "12345678",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/api/v1/user/sign-up", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Use the default HTTP client to send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
