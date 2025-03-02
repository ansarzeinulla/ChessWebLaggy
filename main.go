package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

	log.Fatal(http.ListenAndServe(":8080", nil))
}
