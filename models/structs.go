package models

type Geo struct {
	Results []struct {
		Country string `json:"country"`
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
		Bbox struct {
			Lon1 float64 `json:"lon1"`
			Lat1 float64 `json:"lat1"`
			Lon2 float64 `json:"lon2"`
			Lat2 float64 `json:"lat2"`
		} `json:"bbox"`
	} `json:"results"`
}

type Data struct {
	AllArtists     []Artists
	CurrentArtists []Artists
	Duplicates     map[string]string
	MinDc          int
	MaxDc          int
	HomePage       bool
}

type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	Loca         Locations
	CongertDates string `json:"concertDates"`
	ConDT        Dates
	Relations    string `json:"relations"`
	Rela         Relation
}

type IndexLocations struct {
	Index []Locations `json:"index"`
}

type Locations struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ErrorPage struct {
	Status int
	Type   string
}
