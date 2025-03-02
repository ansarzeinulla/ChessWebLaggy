package main

import (
	"net/http"
	"strings"
)

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

func isValidString(input string) bool {
	return len(input) == 6
}
