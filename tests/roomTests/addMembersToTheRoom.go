package usertests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

var authTokenn string

func setupUserSignupAndSigninAndCreateMember(t *testing.T) {
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
		authTokenn = token
	} else {
		t.Fatalf("Token not found in sign-in response")
	}
}
