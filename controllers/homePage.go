package controllers

import (
	"net/http"

	"groupietracker/models"
	"groupietracker/utils"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderError(w, "Page not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		renderError(w, "Mothod not allowed", http.StatusMethodNotAllowed)
		return
	}

	var artists models.Data

	err := models.FetchAPI("https://groupietrackers.herokuapp.com/api/artists", &artists.AllArtists)
	if err != nil {
		renderError(w, "Server error", http.StatusInternalServerError)
		return
	}

	artists.CurrentArtists = artists.AllArtists

	artists.CurrentArtists[20].Image = "assets/img/3ib.jpg"

	utils.FindMinMax(&artists)

	artists.Duplicates = utils.RemoveDuplicates(artists.AllArtists)

	err = RenderTempalte(w, "./views/index.html", artists, http.StatusOK)
	if err != nil {
		renderError(w, "Server error", http.StatusInternalServerError)
		return
	}
}
