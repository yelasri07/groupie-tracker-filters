package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"groupietracker/models"
	"groupietracker/utils"
)

// HandleFilter fetches data from an API, filters the artists based on user input, and renders the filtered data on the page.
func HandleFilter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		renderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	var artists models.Data

	MINCD1, MAXCD2, FA1, FA2, NOMB, LOC := GetData(r)

	if len(MINCD1) > 100 || len(MAXCD2) > 100 || len(FA1) > 100 ||
		len(FA2) > 100 || len(NOMB) > 100 || len(LOC) > 100 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	filter(&artists, MINCD1, MAXCD2, FA1, FA2, LOC, NOMB)

	artists.Duplicates = utils.RemoveDuplicates(artists.AllArtists)
	artists.HomePage = true

	err := RenderTempalte(w, "views/index.html", artists, http.StatusOK)
	if err != nil {
		renderError(w, "Server error", http.StatusInternalServerError)
		return
	}
}

// artistsFiltred filters the artists based on the provided criteria such as creation date, first album, number of members, and concert locations.
func filter(artists *models.Data, MINCD1, MAXCD2, FA1, FA2, LOC string, NOMB []string) {
	models.FetchAPI("https://groupietrackers.herokuapp.com/api/artists", &artists.AllArtists)
	artists.AllArtists[20].Image = "assets/img/3ib.jpg"

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, artist := range artists.AllArtists {
		wg.Add(1)
		go func(artist models.Artists) {
			defer wg.Done()
			
			hasDate := GetCreattionDate(&artist, MINCD1, MAXCD2)
			hasFirstAlbum := GetFirstAlbum(&artist, FA1, FA2)
			hasMembers := NumberOfMembers(&artist, NOMB)
			hasLocations := LocationsOfConcert(&artist, LOC)

			if hasDate && hasFirstAlbum && hasMembers && hasLocations {
				mu.Lock()
				artists.CurrentArtists = append(artists.CurrentArtists, artist)
				mu.Unlock()
			}
		}(artist)
	}

	wg.Wait()

	utils.FindMinMax(artists)
}

func GetData(r *http.Request) (string, string, string, string, []string, string) {
	MCD1 := r.URL.Query().Get("minCreationDate")
	MCD2 := r.URL.Query().Get("maxCreationDate")
	FA1 := r.URL.Query().Get("firstAlbum1")
	FA2 := r.URL.Query().Get("firstAlbum2")
	NOMB := r.Form["numberOfMembers"]
	LOC := r.URL.Query().Get("locationsOfConcerts")
	return MCD1, MCD2, FA1, FA2, NOMB, LOC
}

func GetCreattionDate(a *models.Artists, min string, max string) bool {
	minV, _ := strconv.Atoi(min)
	maxV, _ := strconv.Atoi(max)
	if minV > maxV {
		minV, maxV = maxV, minV
	}
	if (minV == 1987 && maxV == 1987) || (a.CreationDate >= minV && maxV >= a.CreationDate) {
		return true
	}
	return false
}

func GetFirstAlbum(a *models.Artists, y1, y2 string) bool {
	if len(y1) == 0 && len(y2) == 0 {
		return true
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	minyear, _ := strconv.Atoi(y1)
	maxyear, _ := strconv.Atoi(y2)
	for i := minyear; i <= maxyear; i++ {
		if strings.HasSuffix(a.FirstAlbum, "-"+strconv.Itoa(i)) {
			return true
		}
	}
	return false
}

func LocationsOfConcert(artist *models.Artists, key string) bool {
	var locations models.IndexLocations
	models.FetchAPI("https://groupietrackers.herokuapp.com/api/locations", &locations)

	if key == "" {
		return true
	}
	if key == "seattle-usa" {
		key = "washington-usa"
	}

	idLocation := strings.Split(artist.Locations, "/")[5]
	for _, location := range locations.Index {
		if idLocation == strconv.Itoa(location.ID) {
			for _, loca := range location.Locations {
				if strings.Contains(strings.ToLower(loca), key) {
					return true
				}
			}
		}

		// for _, adress := range locations.Locations {
		// 	if adress == key {
		// 		if locations.ID == artist.ID {
		// 			return true
		// 		}
		// 	}
		// }
	}
	return false
}

func NumberOfMembers(a *models.Artists, key []string) bool {
	if len(key) == 0 {
		return true
	}

	for _, e := range key {
		nb, _ := strconv.Atoi(e)
		if len(a.Members) == nb {
			return true
		}
	}
	return false
}
