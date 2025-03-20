package main

import (
	"encoding/json"
	"os"
)

const (
	firebaseURL     = "https://shahmat-d555c-default-rtdb.firebaseio.com/" // Firebase URL
	credentialsFile = "fire5.json"                                         // Path to your Firebase service account JSON
)

type FirebaseConfig struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
	UniverseDomain          string `json:"universe_domain"`
}

// Function to create firebase.json dynamically
func createFirebaseConfigFile() (string, error) {
	// Define Firebase credentials (DO NOT expose this in public repositories)
	config := FirebaseConfig{
		//NO FIREBASE KEY HERE
	}

	// Convert struct to JSON
	jsonData, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	// Define file path
	filePath := "fire5.json"

	// Write JSON to file
	err = os.WriteFile(filePath, jsonData, 0o644)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
