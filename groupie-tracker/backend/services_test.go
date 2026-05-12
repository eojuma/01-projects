package backend

import (
	"testing"
)

// sampleArtists returns a reusable slice of ArtistDetail for tests.
func sampleArtists() []ArtistDetail {
	return []ArtistDetail{
		{
			Artist: Artist{
				ID:           1,
				Name:         "Queen",
				Members:      []string{"Freddie Mercury", "Brian May", "John Deacon", "Roger Taylor"},
				CreationDate: 1970,
				FirstAlbum:   "14-12-1973",
			},
			DatesLocations: map[string][]string{
				"london-uk": {"01-01-2020"},
			},
		},
		{
			Artist: Artist{
				ID:           2,
				Name:         "Gorillaz",
				Members:      []string{"Damon Albarn", "Jamie Hewlett"},
				CreationDate: 1998,
				FirstAlbum:   "26-03-2001",
			},
			DatesLocations: map[string][]string{},
		},
		{
			Artist: Artist{
				ID:           3,
				Name:         "Foo Fighters",
				Members:      []string{"Dave Grohl", "Nate Mendel", "Taylor Hawkins"},
				CreationDate: 1994,
				FirstAlbum:   "04-07-1995",
			},
			DatesLocations: map[string][]string{
				"new_york-usa": {"15-06-2021"},
			},
		},
	}
}

// SearchArtists tests

func TestSearchArtists_EmptyQuery(t *testing.T) {
	artists := sampleArtists()
	results := SearchArtists("", artists)
	if len(results) != len(artists) {
		t.Errorf("expected %d artists, got %d", len(artists), len(results))
	}
}

func TestSearchArtists_ByName(t *testing.T) {
	artists := sampleArtists()
	results := SearchArtists("queen", artists)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "Queen" {
		t.Errorf("expected Queen, got %s", results[0].Name)
	}
}

func TestSearchArtists_ByMember(t *testing.T) {
	artists := sampleArtists()
	results := SearchArtists("Freddie", artists)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "Queen" {
		t.Errorf("expected Queen, got %s", results[0].Name)
	}
}

func TestSearchArtists_CaseInsensitive(t *testing.T) {
	artists := sampleArtists()
	results := SearchArtists("GORILLAZ", artists)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "Gorillaz" {
		t.Errorf("expected Gorillaz, got %s", results[0].Name)
	}
}

func TestSearchArtists_NoMatch(t *testing.T) {
	artists := sampleArtists()
	results := SearchArtists("nonexistentband", artists)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchArtists_PartialName(t *testing.T) {
	artists := sampleArtists()
	results := SearchArtists("foo", artists)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "Foo Fighters" {
		t.Errorf("expected Foo Fighters, got %s", results[0].Name)
	}
}

func TestSearchArtists_WhitespaceQuery(t *testing.T) {
	artists := sampleArtists()
	results := SearchArtists("  queen  ", artists)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "Queen" {
		t.Errorf("expected Queen, got %s", results[0].Name)
	}
}

func TestSearchArtists_MemberMatchDoesNotDuplicate(t *testing.T) {
	// "Dave Grohl" matches a member of Foo Fighters.
	// Make sure only one copy is returned even though we continue checking members.
	artists := sampleArtists()
	results := SearchArtists("Dave Grohl", artists)
	if len(results) != 1 {
		t.Errorf("expected exactly 1 result, got %d", len(results))
	}
}
