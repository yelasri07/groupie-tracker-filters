package controllers

import (
	"net/http"
	"strconv"

	"groupietracker/models"
)

// info page to view all artists informations
func InfosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	_, err := strconv.Atoi(id)
	if id == "" || err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	artist, err := models.GetArtist("https://groupietrackers.herokuapp.com/api/artists/" + id)
	if err != nil {
		renderError(w, "Server Error", http.StatusInternalServerError)
		return
	}

	if artist.ID == 0 {
		renderError(w, "Artist Not Found", http.StatusNotFound)
		return
	}

	if artist.ID == 21 {
		artist.Image = "./assets/img/3ib.jpg"
	}

	err = RenderTempalte(w, "./views/infos.html", artist, http.StatusOK)
	if err != nil {
		renderError(w, "Server Error", http.StatusInternalServerError)
		return
	}
}
