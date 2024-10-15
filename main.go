package main

import (
	"log"
	"net/http"

	groupie_tracker "groupie_tracker/utils"
)

func main() {
	// Serve specific static files
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		// Define allowed static files
		allowedFiles := map[string]bool{
			"/static/mode.js":       true,
			"/static/error.css":     true,
			"/static/locations.css": true,
			"/static/index.css":     true,
		}

		// Check if the requested file is in the allowed list
		if !allowedFiles[r.URL.Path] {
			groupie_tracker.ErrorHandler(w, http.StatusNotFound)
			return
		}

		// Serve the file if allowed
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	})

	// Set up route handlers
	http.HandleFunc("/", groupie_tracker.ArtistPageHandler)
	http.HandleFunc("/locations/", groupie_tracker.MoreDetailsHandler)
	http.HandleFunc("/favicon.ico", groupie_tracker.FaviconHandler)

	log.Println("Server starting at http://localhost:1234")
	log.Fatal(http.ListenAndServe(":1234", nil))
}
