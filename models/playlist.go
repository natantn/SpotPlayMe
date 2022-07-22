package models

import "gorm.io/gorm"

type Playlist struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	SpotifyID   string `json:"spotify_id"`
	// SpotifyDateCreation string `json:"spotify_date_creation"`
	Musics []Music `gorm:"many2many:playlist_musics;"`
}

type PlaylistMusics struct {
	PlaylistID          int `gorm:"primaryKey"`
	MusicID             int `gorm:"primaryKey"`
	ReprodutionSequence int
}
