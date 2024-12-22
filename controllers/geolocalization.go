package controllers

import (
	"fmt"
	"groupietracker/models"
	"net/http"
)

func GeolocalizationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	location := r.URL.Query().Get("location")

	if location == "" || !IsValidLocation(location) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var Geo models.Geo

	locaAPI := fmt.Sprintf("https://api.geoapify.com/v1/geocode/search?text=%s&format=json&apiKey=5ba3eaa6c92d48c18cccc0c77b034bca", location)
	err := models.FetchAPI(locaAPI, &Geo)
	if err != nil {
		renderError(w, "Server error", http.StatusInternalServerError)
		return
	}

	if len(Geo.Results) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = RenderTempalte(w, "views/geo.html", Geo.Results[0], http.StatusOK)
	if err != nil {
		renderError(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func IsValidLocation(str string) bool {
	var locations models.IndexLocations
	models.FetchAPI("https://groupietrackers.herokuapp.com/api/locations", &locations)

	for _, location := range locations.Index {
		for _, loca := range location.Locations {
			if str == loca {
				return true
			}
		}
	}

	return false
}
