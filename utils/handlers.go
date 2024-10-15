package groupie_tracker

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// function to handle the main page serving
func ArtistPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusNotFound)
		return
	}
	// if r.URL.Path == "/static" {
	// 	println(r.URL.Path)
	// 	// http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	ErrorHandler(w, http.StatusNotFound)
	// 	return
	// }
	artists, err := GetArtists()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	temp, err := template.ParseFiles("template/index.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	temp.Execute(w, artists)
}

// function to display more details about the artist's profile
func MoreDetailsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/locations/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	location, err := GetLocations()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	artist, err := GetArtists()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	relations, err := GetRelations()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	dates, err := GetDates()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	if id > len(artist) || id > len(location) || id > len(relations) || id > len(dates) {
		ErrorHandler(w, http.StatusNotFound)
		return
	}

	data := struct {
		Artist    Artist
		Location  Location
		Date      Date
		Relations Relations
	}{
		Artist:    artist[id-1],
		Location:  location[id-1],
		Date:      dates[id-1],
		Relations: relations[id-1],
	}

	temp, err := template.ParseFiles("template/locations.html")
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		return
	}

	temp.Execute(w, data)
}

// Error handler for displaying errors to users.
func ErrorHandler(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	temp, err := template.ParseFiles("template/error.html")
	if err != nil {
		log.Fatalln(err)
		return
	}
	data := Code{
		Code: code,
	}

	temp.Execute(w, data)
	// temp.Execute(w, map[string]int{"Code": code})
}

// function to handle the favicons in the website
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
