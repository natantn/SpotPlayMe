package models

import (
	"fmt"

	"github.com/natantn/SpotPlayMe/dtos"
	"gorm.io/gorm"
)

type Music struct {
	gorm.Model
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	ReleaseDate string `json:"release_date"`
	SpotifyID   string `json:"spotify_id"`
}

func (m *Music) FillFetchedMusic(mf *dtos.TrackItemInPlaylistResponse) bool {
	updated := false

	if m.Title != mf.Name {
		m.Title = mf.Name
		updated = true
	}
	if m.Album != mf.Album.Name {
		m.Album = mf.Album.Name
		updated = true
	}
	if m.ReleaseDate != mf.Album.ReleaseDate {
		m.ReleaseDate = mf.Album.ReleaseDate
		updated = true
	}

	artist := ""
	for i, artistName := range mf.Artists {
		artist += fmt.Sprintf("%s", artistName.Name)
		if i < len(mf.Artists)-1 {
			artist += ", "
		}
	}

	if m.Artist != artist {
		m.Artist = artist
		updated = true
	}

	return updated
}
