package usertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

var authToken string

// setupUserSignupAndSignin ensures that the user is signed up and signed in before running each test.
func setupUserSignupAndSignin(t *testing.T) {
	// Step 1: Sign up the user
	signupPayload := map[string]string{
		"username": "someone1",
		"email":    "someone1@somewhere.com",
		"password": "12345678",
	}
	signupPayloadBytes, err := json.Marshal(signupPayload)
	if err != nil {
		t.Fatalf("Failed to marshal signup payload: %v", err)
	}

	signupReq, err := http.NewRequest("POST", "http://localhost:8000/api/v1/user/sign-up", bytes.NewBuffer(signupPayloadBytes))
	if err != nil {
		t.Fatalf("Failed to create signup request: %v", err)
	}
	signupReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	signupResp, err := client.Do(signupReq)
	if err != nil || signupResp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to sign up: %v", err)
	}
	defer signupResp.Body.Close()

	// Step 2: Sign in the user to get the authentication token
	signinPayload := map[string]string{
		"username": "someone1@somewhere.com",
		"password": "12345678",
	}
	signinPayloadBytes, err := json.Marshal(signinPayload)
	if err != nil {
		t.Fatalf("Failed to marshal signin payload: %v", err)
	}

	signinReq, err := http.NewRequest("POST", "http://localhost:8000/api/v1/user/sign-in", bytes.NewBuffer(signinPayloadBytes))
	if err != nil {
		t.Fatalf("Failed to create signin request: %v", err)
	}
	signinReq.Header.Set("Content-Type", "application/json")

	signinResp, err := client.Do(signinReq)
	if err != nil || signinResp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to sign in: %v", err)
	}
	defer signinResp.Body.Close()

	var signinRespBody map[string]interface{}
	if err := json.NewDecoder(signinResp.Body).Decode(&signinRespBody); err != nil {
		t.Fatalf("Failed to decode signin response: %v", err)
	}

	// Extract the auth token from the response and store it in the global variable
	if token, ok := signinRespBody["access-token"].(string); ok {
		authToken = token
	} else {
		t.Fatalf("Token not found in sign-in response")
	}
}

func TestCreateRoom(t *testing.T) {
	setupUserSignupAndSignin(t)

	createRoomPayload := map[string]string{
		"name": "someroom2someone1",
	}
	createRoomPayloadBytes, err := json.Marshal(createRoomPayload)
	req, err := http.NewRequest("POST", "http://localhost:8000/api/v1/room/create-room", bytes.NewBuffer(createRoomPayloadBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if err != nil {
		t.Fatalf("Failed to marshal signin payload: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code and body as needed
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
