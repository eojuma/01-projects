package backend

import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
	"sync"
)

// cache stores fetched data to avoid repeated API calls
var (
	cachedArtists []ArtistDetail
	cacheMu       sync.Mutex
	cacheLoaded   bool
)

func GetArtists() ([]Artist, error) {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}


func GetLocations() ([]Location, error) {
	url := "https://groupietrackers.herokuapp.com/api/locations"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var locations []Location
	err = json.NewDecoder(resp.Body).Decode(&locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}



func GetDates() ([]Date, error) {
	url := "https://groupietrackers.herokuapp.com/api/dates"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dates []Date
	err = json.NewDecoder(resp.Body).Decode(&dates)
	if err != nil {
		return nil, err
	}

	return dates, nil
}

func GetRelations() (map[int]map[string][]string, error) {
	url := "https://groupietrackers.herokuapp.com/api/relation"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch relations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("relations API returned status %d", resp.StatusCode)
	}

	var relIndex RelationIndex
	if err := json.NewDecoder(resp.Body).Decode(&relIndex); err != nil {
		return nil, fmt.Errorf("failed to decode relations: %w", err)
	}

	result := make(map[int]map[string][]string)
	for _, r := range relIndex.Index {
		result[r.ID] = r.DatesLocations
	}
	return result, nil
}

func GetArtistsWithRelations() ([]ArtistDetail, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if cacheLoaded {
		return cachedArtists, nil
	}

	artists, err := GetArtists()
	if err != nil {
		return nil, err
	}

	relations, err := GetRelations()
	if err != nil {
		return nil, err
	}

	details := make([]ArtistDetail, len(artists))
	for i, a := range artists {
		dl := relations[a.ID]
		if dl == nil {
			dl = make(map[string][]string)
		}
		details[i] = ArtistDetail{
			Artist:         a,
			DatesLocations: dl,
		}
	}

	cachedArtists = details
	cacheLoaded = true
	return cachedArtists, nil
}

func SearchArtists(query string, artists []ArtistDetail) []ArtistDetail {
	if query == "" {
		return artists
	}
	q := strings.ToLower(strings.TrimSpace(query))
	var results []ArtistDetail
	for _, a := range artists {
		if strings.Contains(strings.ToLower(a.Name), q) {
			results = append(results, a)
			continue
		}
		for _, m := range a.Members {
			if strings.Contains(strings.ToLower(m), q) {
				results = append(results, a)
				break
			}
		}
	}
	return results
}