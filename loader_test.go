// main_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	groupie_tracker "groupie_tracker/utils"
)

func TestArtistPageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(groupie_tracker.ArtistPageHandler)

	handler.ServeHTTP(rr, req)

	// Check if the response code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Add more checks here to verify the response body if needed
}

func TestMoreDetailsHandler_ValidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/locations/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(groupie_tracker.MoreDetailsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Add more checks here to verify the response body if needed
}

func TestMoreDetailsHandler_InvalidID(t *testing.T) {
	req, err := http.NewRequest("GET", "/locations/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(groupie_tracker.MoreDetailsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestGetArtists(t *testing.T) {
	artists, err := groupie_tracker.GetArtists()
	if err != nil {
		t.Fatalf("GetArtists failed: %v", err)
	}

	if len(artists) == 0 {
		t.Error("expected non-empty artist list")
	}
}

func TestGetLocations(t *testing.T) {
	locations, err := groupie_tracker.GetLocations()
	if err != nil {
		t.Fatalf("GetLocations failed: %v", err)
	}

	if len(locations) == 0 {
		t.Error("expected non-empty location list")
	}
}

func TestGetRelations(t *testing.T) {
	relations, err := groupie_tracker.GetRelations()
	if err != nil {
		t.Fatalf("GetRelations failed: %v", err)
	}

	if len(relations) == 0 {
		t.Error("expected non-empty relations list")
	}
}

func TestGetDates(t *testing.T) {
	dates, err := groupie_tracker.GetDates()
	if err != nil {
		t.Fatalf("GetDates failed: %v", err)
	}

	if len(dates) == 0 {
		t.Error("expected non-empty dates list")
	}
}
