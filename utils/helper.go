package utils

import (
	"strconv"

	"groupietracker/models"
)

// Remove duplicates from input search like "Locations"
func RemoveDuplicates(artists []models.Artists) map[string]string {
	var locations models.IndexLocations

	models.FetchAPI("https://groupietrackers.herokuapp.com/api/locations", &locations)

	duplicates := make(map[string]string)
	for _, artist := range artists {
		duplicates[artist.Name] = "artist/band"
		duplicates[artist.FirstAlbum] = "first album"
		duplicates[strconv.Itoa(artist.CreationDate)] = "creation date"
		for _, member := range artist.Members {
			duplicates[member] = "member"
		}
	}

	for _, location := range locations.Index {
		for _, loca := range location.Locations {
			duplicates[loca] = "location"
		}
	}

	return duplicates
}


// find min and max value creation date
func FindMinMax(artists *models.Data) {
	min := 2024
	max := 0
	for _, ele := range artists.AllArtists {
		if min > ele.CreationDate {
			min = ele.CreationDate
		} else if max < ele.CreationDate {
			max = ele.CreationDate
		}
	}
	artists.MaxDc = max
	artists.MinDc = min
}

// replace "-" to ", "
func Replace(location string) string {
	if len(location) == 0  {
		return ""
	}

	var str string
	for _, char := range location {
		if char == '-' {
			str += ", "
			continue
		}
		str += string(char)
	}
	return str
}