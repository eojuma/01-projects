package backend

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Location struct {
	ID        int      `json:"id"`
	Cities    []string `json:"cities"`
	Countries []string `json:"countries"`
}




type Date struct {
	ID      int     `json:"id"`    
	Dates   []string `json:"dates"` 
}




type Relation struct {
	ArtistID   int   `json:"artistId"`
	LocationID int   `json:"locationId"`
	Dates      []int `json:"dates"` 
}

type ArtistDetail struct {
	Artist
	DatesLocations map[string][]string
}

type RelationResponse struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// PageData holds data passed to the home page template
type PageData struct {
	Artists []ArtistDetail
	Query   string
	Total   int
}

type ErrorData struct {
	StatusCode int
	Message    string
}

type RelationIndex struct {
	Index []RelationResponse `json:"index"`
}