package main

import (
	"bytes"
	"context"
	"fmt"
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
	go startGameLoop()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func startGameLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			// Simulate a JSON request body for POST
			jsonData := []byte(`{"from":"e2","to":"e4","fen":""}`)
			resp, err := http.Post("http://localhost:8080/game", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Printf("Error making HTTP request: %v", err)
				return
			}
			resp.Body.Close() // Close response to prevent leaks
		}()
	}
}

type fakeResponseWriter struct{}

func (f *fakeResponseWriter) Header() http.Header        { return http.Header{} }
func (f *fakeResponseWriter) Write([]byte) (int, error)  { return 0, nil }
func (f *fakeResponseWriter) WriteHeader(statusCode int) {}
