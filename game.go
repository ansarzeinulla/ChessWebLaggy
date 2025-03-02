package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/notnil/chess"
)

func gameHandler(w http.ResponseWriter, r *http.Request) {
	// Handle GET request (to load initial game page)
	if r.Method == http.MethodGet {
		// Retrieve the current game FEN from Firebase (or initialize it)
		ctx := context.Background()
		dbClient, err := client.Database(ctx)
		if err != nil {
			log.Printf("Error getting Firebase database client: %v", err)
			http.Error(w, "Error connecting to Firebase Database", http.StatusInternalServerError)
			return
		}

		// Firebase path reference for the game FEN (replace `gameID` with actual game identifier)
		firebaseRef := dbClient.NewRef("games-" + gameID + "-aboba")

		// Get the current FEN from Firebase
		var currentFEN string
		if err := firebaseRef.Get(ctx, &currentFEN); err != nil {
			log.Printf("Error retrieving FEN from Firebase: %v", err)
			http.Error(w, "Error retrieving FEN from Firebase", http.StatusInternalServerError)
			return
		}

		// If there's no FEN, initialize a new game
		if currentFEN == "" {
			// Set the initial FEN for a new game
			currentFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
			if err := firebaseRef.Set(ctx, currentFEN); err != nil {
				log.Printf("Error setting initial FEN to Firebase: %v", err)
				http.Error(w, "Error setting initial FEN to Firebase", http.StatusInternalServerError)
				return
			}
		}

		// Render the game page (HTML template)
		data := struct {
			GameID      string
			PlayerColor bool
			FEN         string
		}{
			GameID:      gameID,
			PlayerColor: playerColor,
			FEN:         currentFEN,
		}

		// Parse and execute the template to send the game page with FEN data
		tmpl, err := template.ParseFiles("templates/game.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Handle POST request (to handle move updates)
	if r.Method == http.MethodPost {
		var moveData struct {
			From string `json:"from"`
			To   string `json:"to"`
			FEN  string `json:"fen"`
		}

		// Parse the JSON body
		err := json.NewDecoder(r.Body).Decode(&moveData)
		if err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		// Retrieve or create the game
		ctx := context.Background()
		dbClient, err := client.Database(ctx)
		if err != nil {
			log.Printf("Error getting Firebase database client: %v", err)
			http.Error(w, "Error connecting to Firebase Database", http.StatusInternalServerError)
			return
		}

		// Firebase path reference for the game FEN
		firebaseRef := dbClient.NewRef("games-" + gameID + "-aboba")

		// Get the current FEN from Firebase
		var currentFEN string
		if err := firebaseRef.Get(ctx, &currentFEN); err != nil {
			log.Printf("Error retrieving FEN from Firebase: %v", err)
			http.Error(w, "Error retrieving FEN from Firebase", http.StatusInternalServerError)
			return
		}

		// If no FEN is found, initialize a new game
		if currentFEN == "" {
			currentFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
		}

		// Create a new chess game from the current FEN
		fenStr := currentFEN
		fen, err := chess.FEN(fenStr)
		if err != nil {
			log.Printf("Error parsing FEN: %v", err)
			http.Error(w, "Invalid FEN string", http.StatusBadRequest)
			return
		}

		game := chess.NewGame(chess.UseNotation(chess.UCINotation{}), fen)

		// Apply the move using UCI notation
		move := fmt.Sprintf("%s%s", moveData.From, moveData.To)
		if err := game.MoveStr(move); err != nil {
			log.Printf("Error applying move: %v", err)
			http.Error(w, "Invalid move", http.StatusBadRequest)
			return
		}

		// Get the new FEN after the move
		newFEN := game.Position().String()

		// Update Firebase with the new FEN
		if err := firebaseRef.Set(ctx, newFEN); err != nil {
			log.Printf("Error updating Firebase: %v", err)
			http.Error(w, "Error updating Firebase", http.StatusInternalServerError)
			return
		}

		// Send the new FEN and move validity response back
		response := struct {
			Valid bool   `json:"valid"`
			FEN   string `json:"fen"`
		}{
			Valid: true,
			FEN:   newFEN,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func fetchFENHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	dbClient, err := client.Database(ctx)
	if err != nil {
		http.Error(w, "Error connecting to Firebase Database", http.StatusInternalServerError)
		return
	}

	firebaseRef := dbClient.NewRef("games-" + gameID + "-aboba")
	var currentFEN string
	if err := firebaseRef.Get(ctx, &currentFEN); err != nil {
		http.Error(w, "Error retrieving FEN from Firebase", http.StatusInternalServerError)
		return
	}
	if currentFEN == "" {
		currentFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	}

	response := struct {
		FEN string `json:"fen"`
	}{
		FEN: currentFEN,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func serveGamePageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/game.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		GameID      string
		PlayerColor bool
	}{
		GameID:      gameID,
		PlayerColor: playerColor,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
