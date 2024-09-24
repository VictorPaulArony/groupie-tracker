package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Location     string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relation     string   `json:"relations"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	// Dates     string   `json:"dates"`
}
type LocationIndex struct {
	Index []Location `json:"index"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DateIndex struct {
	Index []Date `json:"index"`
}

type Relations struct {
	ID       int                 `json:"id"`
	Relation map[string][]string `json:"datesLocations"`
}

type RelationIndex struct {
	Index []Relations `json:"index"`
}

// function to get the APIs and unmershal them
func GetApis(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("ERROR MAKING GET REQUEST: ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("ERROR READING RESPONSE BODY: ", err)
	}

	return json.Unmarshal(body, target)
}

// function to get the Artist
func GetArtists() ([]Artist, error) {
	var artists []Artist
	err := GetApis("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

// function to get the Locations to visit
func GetLocations() ([]Location, error) {
	var index LocationIndex

	err := GetApis("https://groupietrackers.herokuapp.com/api/locations", &index)
	if err != nil {
		return nil, err
	}
	return index.Index, nil
}

// function to unmarshal the Relations
func GetRelations() ([]Relations, error) {
	var relationIndex RelationIndex

	err := GetApis("https://groupietrackers.herokuapp.com/api/relation", &relationIndex)
	if err != nil {
		return nil, err
	}

	return relationIndex.Index, nil
}

// function to get all the dates in the API
func GetDates() ([]Date, error) {
	var dateindex DateIndex

	err := GetApis("https://groupietrackers.herokuapp.com/api/dates", &dateindex)
	if err != nil {
		return nil, err
	}

	return dateindex.Index, nil
}

// function to haldle the main page serving
func ArtistPageHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := GetArtists()
	if err != nil {
		http.Error(w, "Unable To Fatch Artist Data", http.StatusInternalServerError)
	}

	temp, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Redirect(w, r, "/404.html", http.StatusNotFound)
	}

	temp.Execute(w, artists)
}

// function to display more details abount the artists profile
func MoreDetailsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/locations/")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > 52 {
		http.Redirect(w, r, "/404.html", http.StatusNotFound)
	}
	location, err := GetLocations()
	if err != nil {
		log.Fatalln("UNNABLE TO FETCH LOCATION DATA: ", err)
	}
	artist, err := GetArtists()
	if err != nil {
		log.Fatalln("UNNABLE TO FETCH ARTIST DATA: ", err)
	}
	relations, err := GetRelations()
	if err != nil {
		log.Fatalln("UNNABLE TO FETCH RELATIONS DATA: ", err)
	}
	dates, err := GetDates()
	if err != nil {
		log.Fatalln("UNNABLE TO FETCH DATES DATA: ", err)
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
		http.Redirect(w, r, "/404.html", http.StatusNotFound)
	}

	temp.Execute(w, data)
}

func main() {
	http.HandleFunc("/", ArtistPageHandler)
	http.HandleFunc("/locations/", MoreDetailsHandler)
	log.Println("http://localhost:1234")

	log.Fatal(http.ListenAndServe(":1234", nil))
}
