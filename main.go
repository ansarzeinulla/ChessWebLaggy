package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/notnil/chess"
	"google.golang.org/api/option"
)

// Global variables to store game details
var (
	gameID      string
	playerColor bool
	client      *firebase.App
)

const (
	firebaseURL     = "https://shahmat-d555c-default-rtdb.firebaseio.com/" // Firebase URL
	credentialsFile = "fire.json"                                          // Path to your Firebase service account JSON
)

func main() {
	// Serve the static HTML pages
	ctx := context.Background()
	opt := option.WithCredentialsFile(credentialsFile)
	config := &firebase.Config{DatabaseURL: firebaseURL}

	var err error
	client, err = firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing Firebase app: %v", err)
	}
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/game", gameHandler)
	http.Handle("/figures/", http.StripPrefix("/figures/", http.FileServer(http.Dir("figures"))))

	// Start the web server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handler to serve the index.html page and process the POST request
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Handle form submission
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Retrieve the input string and color
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

// Game handler to serve the game page
func gameHandler(w http.ResponseWriter, r *http.Request) {
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

		// Update Firebase with the new game state (move)
		dbClient, err := client.Database(context.Background())
		if err != nil {
			log.Printf("Error connecting to Firebase Database: %v", err)
			http.Error(w, "Error connecting to Firebase Database", http.StatusInternalServerError)
			return
		}

		// Corrected Firebase path reference
		firebaseRef := dbClient.NewRef("games/" + moveData.FEN + "/aboba") // Ensure this matches your Firebase structure

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
		GameID:      gameID,                                                     // You can replace this with dynamic game ID
		PlayerColor: playerColor,                                                // Or dynamically set based on user
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
