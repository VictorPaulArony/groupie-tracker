package groupie_tracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
type Code struct {
	Code int
}

// function to get the APIs and unmarshal them
func GetApis(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
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
