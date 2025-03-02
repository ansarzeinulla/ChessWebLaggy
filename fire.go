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
		Type:                    "service_account",
		ProjectID:               "shahmat-d555c",
		PrivateKeyID:            "a446b652f8fe1394f8a42cdf7f35914c9fda41f9",
		PrivateKey:              "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDkeDBr5LNAcehw\n0wOxymsUBjXNh0TBAJ4i00VAAT4fE1UVLdxGpXA2yPAVw52y4oQylp9WRMi2Tarf\nxEFuIV0XnEbU77KxGyFl/Ura3ba782CdbL5ReQJzwap4fk7HjrLIUa0hhxjtYHPU\nutBOUlg5lLVDCk/RHqiWBKL2BY+CxpSS4/KfkwTUJ6FF1qaNFux7ZEXxMt3eygki\nFxdoTfBqJE9wgaPqyYdvils6QBbWNafj7no/b9rqerr0wzwqjNd6cdXYOO6w7KdO\n+zMmsFA8+tvxN42TBqKSRdXcCJVAIba9Ou9cdDhtnJy+BuxaBDrvq4t1ugKUYnU+\ncaWVvlEFAgMBAAECggEACBUibWa5WrHM7dThQBRvplUz6RH8gjFx5rx/qiyBFJtb\n+rRVIyCuMBn6uGJjvVTbJWkDuQsU6LIlteXdBLbkcrXiIZ1bPn29uDx2PfjzAoIx\nLFCRCRPVXOAAlmOtiMzJEWImVFXRsXAkmOWFAghSu7LBmN4QBc9mWmxOkZKPyaoI\nrzf2npoKHnJX/pLKD4Di+3mSwKgK6sDOLhtH9kaUILVnL0rq4JfflgGXVXrFb3Jm\n/qIj0rsuCcnitlISSWPIBTctvqPzzVDaHq+CE2yMBV4tclHgJTteCO1OSpBntjAO\nwTIJ9jXvGgh6tY5+/5h/ocU9p9yARwbsY5i6isUo3wKBgQD6kT0pgNAHl8wt6doW\nE0cfdUegxqNVtNTxFaxjEFYipi8LEXztqyNvgLI52vB8P4ncGmfnGIhWZBgrPCao\nnseUAskfb5WIUSJWkhvzcbj2ZsA0IweGrX7y6T+rVCrN/UbaVTYYMvGReTuBnrwx\nHU86JE+nAkKGmkPXCw6pEQTUJwKBgQDpbEwYuRsh1xwUBilvCueKu/YVtdFxOEl1\nlKy5zY6IVk+1sRoA44iz9uKVFKhxbzcXscV2LpP1Q20NBVHQDvABKa7mdbNqppr0\nNHBOkLxzL89Z9lfyPKWifPxQo4D/kCv30RpMj+7TWHZx2abeHlfJZBQhR3LU2+YJ\n9vna13WQ8wKBgEc6Xk6cBYcDCdHLdmlsFX3F0xTLIsdMXnQiGx0WGcZDw3+7+u19\nBte9l+yGZnKLhV8CSqMRAEC+t3gi40Jv0IAswoujJrjXh5Fge32ayF+TGfQ4OP15\n+GqJD8ZeaMShyTBrpLMAWFdoRRg1zX2QvWLjy5jINa0Z0UsiI4rAAcVlAoGACWcj\naZuLTEGuD+Bvqtl1mlEYCKfaWAU8cFAc5R8yrqtLarZHpeGEkDtRxU+fuXIRdhLj\nMW+O5kJhEjU0pnzzjhhvwzjakWFEvLGgFIogDUPPxn/16vwmb/U49MahW6ojG0iB\nFrR1mm3l15A8+JWgU6yEYxLNvWVeTuh0CCzFv6ECgYAHbJC9PnDWfe1snOkvKzV0\nSuTkoq/rEU6YvFGjFbjVbej5YKjEZmtyTNcqM5tI3yaedYRGlQdJUBsV9xjpqf84\nTn5oq+nEiEdRU04PWSBo9L2ZrDADlCnXkeAh0LJZZRv6fGKlqk72HG90x8N1oTA3\nkSkXWsZxRAcDJGUe06w3wQ==\n-----END PRIVATE KEY-----\n",
		ClientEmail:             "firebase-adminsdk-fbsvc@shahmat-d555c.iam.gserviceaccount.com",
		ClientID:                "100510951141045243794",
		AuthURI:                 "https://accounts.google.com/o/oauth2/auth",
		TokenURI:                "https://oauth2.googleapis.com/token",
		AuthProviderX509CertURL: "https://www.googleapis.com/oauth2/v1/certs",
		ClientX509CertURL:       "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-fbsvc%40shahmat-d555c.iam.gserviceaccount.com",
		UniverseDomain:          "googleapis.com",
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
