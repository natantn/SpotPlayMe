package models

import "gorm.io/gorm"

type Music struct {
	gorm.Model
	Title string `json:"title"`
	Artist string `json:"artist"`
	Album string `json:"album"`
	ReleaseDate string `json:"release_date"`
}