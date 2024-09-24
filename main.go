/*package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*************************************************************************************/
//structs for data manipulation during unmarshaling of the json API
/*
type ApiIndex struct {
	Artist   string `json:"artists"`
	Location string `json:"locations"`
	Date     string `json:"dates"`
	Relation string `json:"relation"`
}

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

/*******************************************************************************************/
//the unmarshaling of the APIs data using get method
/*
func GetApis() ApiIndex {
	url := "https://groupietrackers.herokuapp.com/api"
	res, err := http.Get(url)
	if err != nil {
		log.Println("ERROR MAKING GET REQUEST:", err)
		return ApiIndex{}
	}
	defer res.Body.Close()

	var api ApiIndex
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("ERROR READING RESPONSE BODY:", err)
		return ApiIndex{}
	}
	err = json.Unmarshal(body, &api)
	if err != nil {
		log.Println("ERROR UNMARSHALLING RESPONSE BODY:", err)
		return ApiIndex{}
	}
	return api
}

func GetArtists(url string) ([]Artist, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func GetLocations(url string) (Location, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Location{}, err
	}
	var location Location
	err = json.Unmarshal(body, &location)
	if err != nil {
		return Location{}, err
	}
	return location, nil
}

func GetDates(url string) (Date, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Date{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Date{}, err
	}
	var date Date
	err = json.Unmarshal(body, &date)
	if err != nil {
		log.Println("Error Unmarshalling Dates:", err)
		return Date{}, err
	}

	return date, nil
}

func GetRelations(url string) (Relation, error) {
	resp, err := http.Get(url)
	if err != nil {
		return Relation{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Relation{}, err
	}
	var relation Relation
	err = json.Unmarshal(body, &relation)
	if err != nil {
		return Relation{}, err
	}
	return relation, nil
}

func ArtistPageHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := GetArtists(GetApis().Artist)
	if err != nil {
		http.Error(w, "UNABLE TO FETCH ARTIST DATA", http.StatusInternalServerError)
		return
	}

	if len(artists) == 0 {
		http.Error(w, "NO ARTISTS FOUND", http.StatusNotFound)
		return
	}

	templ, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Fatalln("Error While Parsing The file")
	}

	templ.Execute(w, artists)
}

func FormatLocation(location string) string {
	location = strings.ReplaceAll(location, "-", ", ")
	location = strings.ReplaceAll(location, "_", " ")

	words := strings.Fields(location)

	for i, word := range words {
		words[i] = strings.ToTitle(strings.ToLower(word))
	}

	return strings.Join(words, " ")
}

/************************************************************************************************/
//http/net handlers for the post method of the data and parsing the html files
/*
func LocationPageHandler(w http.ResponseWriter, r *http.Request) {
	strId := strings.TrimPrefix(r.URL.Path, "/locations/")
	id, err := strconv.Atoi(strId)
	if err != nil || id < 1 || id > 52 {
		http.Error(w, "INVALID LOCATION ID", http.StatusBadRequest)
		return
	}

	location, err := GetLocations(GetApis().Location + "/" + strconv.Itoa(id))
	if err != nil {
		http.Error(w, "UNNABLE TO FETCH LOCATION DATA", http.StatusInternalServerError)
		return
	}

	relations, err := GetRelations(GetApis().Relation + "/" + strconv.Itoa(id))
	if err != nil {
		http.Error(w, "UNNABLE TO FETCH RELATION DATA", http.StatusInternalServerError)
		return
	}

	date, err := GetDates(location.Dates)
	if err != nil {
		http.Error(w, "UNNABLE TO FETCH DATES DATA", http.StatusInternalServerError)
		return
	}

	artists, err := GetArtists(GetApis().Artist)
	if err != nil {
		http.Error(w, "UNNABLE TO FETCH ARTIST DATA", http.StatusInternalServerError)
		return
	}

	if id > len(artists) {
		http.Error(w, "INVALID LOCATION ID", http.StatusBadRequest)
		return
	}

	formatLocation := make(map[string][]string)
	for location, datesList := range relations.DatesLocations {
		location = FormatLocation(location)
		formatLocation[location] = datesList
	}

	data := struct {
		Artist   Artist
		Location Location
		Relation Relation
		Dates    Date
	}{
		Artist:   artists[id-1],
		Location: location,
		Relation: Relation{DatesLocations: formatLocation},
		Dates:    date,
	}

	templ, err := template.ParseFiles("template/locations.html")
	if err != nil {
		log.Fatalln("Error While Parsing The file")
	}

	templ.Execute(w, data)
}

/************************************************************************************************/
//the main function to sarve and listen to the http/net requests
/*
func main() {
	file := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", file))

	http.HandleFunc("/", ArtistPageHandler)
	http.HandleFunc("/locations/", LocationPageHandler)
	log.Println("Server is running on port http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// end of code

*/

