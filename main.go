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

var (
	gameID      string
	playerColor bool
	client      *firebase.App
)

const (
	firebaseURL     = "https://shahmat-d555c-default-rtdb.firebaseio.com/"    // Firebase URL
	credentialsFile = "shahmat-d555c-firebase-adminsdk-fbsvc-b543ad59de.json" // Path to your Firebase service account JSON
)

func main() {
	ctx := context.Background()
	opt := option.WithCredentialsFile(credentialsFile)
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
