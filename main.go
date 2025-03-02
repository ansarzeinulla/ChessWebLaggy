package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var (
	gameID      string
	playerColor bool
	client      *firebase.App
)

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
	client = app
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/game", gameHandler)
	log.Println("Server started on :8080")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/figures/", http.StripPrefix("/figures/", http.FileServer(http.Dir("figures"))))
	http.HandleFunc("/fetch-fen", fetchFENHandler)
	go startGameLoop()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startGameLoop() {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			// 🔹 Fetch the latest FEN before making a move
			resp, err := http.Get("http://localhost:8080/fetch-fen")
			if err != nil {
				log.Printf("Error fetching FEN: %v", err)
				return
			}
			defer resp.Body.Close()

			var fenData struct {
				FEN string `json:"fen"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&fenData); err != nil {
				log.Printf("Error decoding FEN response: %v", err)
				return
			}

			// 🔹 If FEN hasn't changed, don't send a move
			if fenData.FEN == "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" {
				jsonData := `{"from":"e2","to":"e4","fen":""}`
				resp, err := http.Post("http://localhost:8080/game", "application/json", bytes.NewBuffer([]byte(jsonData)))
				if err != nil {
					log.Printf("Error making HTTP request: %v", err)
					return
				}
				defer resp.Body.Close()

				body, _ := io.ReadAll(resp.Body)
				log.Printf("Move Response: %s", body)
			}
		}()
	}
}

type fakeResponseWriter struct {
	header http.Header
	body   bytes.Buffer
	status int
}

func (f *fakeResponseWriter) Header() http.Header {
	if f.header == nil {
		f.header = make(http.Header)
	}
	return f.header
}

func (f *fakeResponseWriter) Write(b []byte) (int, error) {
	return f.body.Write(b)
}

func (f *fakeResponseWriter) WriteHeader(statusCode int) {
	f.status = statusCode
}
