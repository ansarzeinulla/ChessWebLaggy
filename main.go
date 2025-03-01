package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/notnil/chess"
	"google.golang.org/api/option"
)

var (
	gameID      string
	playerColor bool
	client      *firebase.App
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

func main() {
	filePath, err := createFirebaseConfigFile()
	if err != nil {
		fmt.Println("Error creating Firebase config file:", err)
		return
	}

	ctx := context.Background()
	opt := option.WithCredentialsFile(filePath)
	config := &firebase.Config{DatabaseURL: firebaseURL}

	// Create Firebase app
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v", err)
		return
	}
	client = app // Store app for later use

	// Set up HTTP handlers
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/game", gameHandler)
	http.Handle("/figures/", http.StripPrefix("/figures/", http.FileServer(http.Dir("figures"))))

	// Start the server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Handle form submission
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Retrieve input string and color
		inputString := r.FormValue("inputString")
		colorChecked := r.FormValue("color")

		// Validate the input string
		if len(inputString) != 6 || !isValidString(inputString) {
			// If invalid, stay on the index page
			http.ServeFile(w, r, "templates/index.html")
			return
		}

		// Capitalize the input string
		inputString = strings.ToUpper(inputString)

		// Store data in global variables
		gameID = inputString
		playerColor = colorChecked == "on"

		// Redirect to the game page
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}

	// Serve the index.html page for GET requests
	http.ServeFile(w, r, "templates/index.html")
}

// Helper function to validate the string (only English letters)
func isValidString(input string) bool {
	re := regexp.MustCompile(`^[A-Za-z]+$`)
	return re.MatchString(input)
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	// Get the Firebase Realtime Database client
	ctx := context.Background()
	dbClient, err := client.Database(ctx)
	if err != nil {
		log.Printf("Error getting Firebase database client: %v", err)
		http.Error(w, "Error connecting to Firebase Database", http.StatusInternalServerError)
		return
	}

	// Handle POST request
	if r.Method == "POST" {
		var moveData struct {
			From string `json:"from"`
			To   string `json:"to"`
			FEN  string `json:"fen"`
		}

		// Decode the move data from the request body
		err := json.NewDecoder(r.Body).Decode(&moveData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Parse the FEN string
		officialfen, err := chess.FEN(moveData.FEN)
		if err != nil {
			http.Error(w, "Invalid FEN string", http.StatusBadRequest)
			return
		}

		// Initialize a new chess game with the given FEN
		game := chess.NewGame(chess.UseNotation(chess.UCINotation{}), officialfen)

		// Create move string and validate the move format
		moveStr := moveData.From + moveData.To
		if len(moveStr) != 4 {
			http.Error(w, "Invalid move format", http.StatusBadRequest)
			return
		}

		// Try making the move
		if err := game.MoveStr(moveStr); err != nil {
			http.Error(w, "Invalid move", http.StatusBadRequest)
			return
		}

		// Firebase path reference
		firebaseRef := dbClient.NewRef("games-" + gameID + "-aboba")

		// Set the move in Firebase
		err = firebaseRef.Set(context.Background(), moveStr)
		if err != nil {
			log.Printf("Error updating Firebase: %v", err)
			http.Error(w, "Error updating Firebase", http.StatusInternalServerError)
			return
		}

		// Respond with the updated FEN and valid move status
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": true,
			"fen":   game.Position().String(),
		})
		return
	}

	// If not POST, render the game page
	data := struct {
		GameID      string
		PlayerColor bool
		FEN         string
	}{
		GameID:      gameID,
		PlayerColor: playerColor,
		FEN:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", // Default starting FEN
	}

	// Parse the template and execute it
	tmpl, err := template.ParseFiles("templates/game.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
