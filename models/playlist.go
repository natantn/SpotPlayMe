package models

import (
	spotify "github.com/natantn/SpotPlayMe/integrations/spotify"
	"gorm.io/gorm"
)

type Playlist struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	SpotifyID   string `json:"spotify_id"`
	// SpotifyDateCreation string `json:"spotify_date_creation"`
	Musics []Music `gorm:"many2many:playlist_musics;"`
}

func (p *Playlist) FillFetchedPlaylist(pf *spotify.PlaylistApiResponse, musics *[]Music) (hasChanged bool) {
	hasChanged = false
	if p.Title != pf.Name {
		p.Title = pf.Name
		hasChanged = true
	}
	if p.Description != pf.Description {
		p.Description = pf.Description
		hasChanged = true
	}
	if p.SpotifyID != pf.ID {
		p.SpotifyID = pf.ID
		hasChanged = true
	}
	p.Musics = *musics

	return
}

type PlaylistMusics struct {
	PlaylistID          int `gorm:"primaryKey"`
	MusicID             int `gorm:"primaryKey"`
	ReprodutionSequence int
}

func (pm *PlaylistMusics) SetMusicSequenceInPlaylist(seq int) {
	pm.ReprodutionSequence = seq
}
