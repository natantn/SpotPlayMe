package services

import (
	"errors"

	"github.com/natantn/SpotPlayMe/dtos"
	database "github.com/natantn/SpotPlayMe/integrations/database"
	"github.com/natantn/SpotPlayMe/models"
	"gorm.io/gorm"
)

func GetMusicsByTitle(title string) *[]models.Music {
	musics := []models.Music{}
	database.DB.Where("title LIKE ?", title).Find(&musics)

	return &musics
}

func GetMusicOccorrencesInPlaylists(musics *[]models.Music) *[]dtos.MusicFound {
	musicsFound := []dtos.MusicFound{}

	for _, music := range *musics {
		occurrences := []models.PlaylistMusics{}
		database.DB.Where(&models.PlaylistMusics{MusicID: int(music.ID)}).Find(&occurrences)

		musicFound := dtos.MusicFound{
			Music: dtos.MusicItem{
				Title:  music.Title,
				Artist: music.Artist,
				Album:  music.Album,
			},
		}

		for _, occurrence := range occurrences {
			playlist := models.Playlist{}
			database.DB.First(&playlist, occurrence.PlaylistID)

			seq := occurrence.ReprodutionSequence

			findMusicOccorrunceInPlaylist := func(playlistID, seq int) (music models.Music) {
				occorruces := models.PlaylistMusics{}
				err := database.DB.Where("playlist_id = ? AND reprodution_sequence = ?", playlistID, seq).First(&occorruces).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return models.Music{}
				} else {
					database.DB.First(&music, occorruces.MusicID)
					return
				}

			}

			previousMusic := findMusicOccorrunceInPlaylist(occurrence.PlaylistID, seq-1)
			nextMusic := findMusicOccorrunceInPlaylist(occurrence.PlaylistID, seq+1)

			musicOccorruce := dtos.MusicOccorruce{
				PlaylistTitle:        playlist.Title,
				ReproduntionSequence: occurrence.ReprodutionSequence,
				Previous: dtos.MusicItem{
					Title:  previousMusic.Title,
					Artist: previousMusic.Artist,
					Album:  previousMusic.Album,
				},
				Next: dtos.MusicItem{
					Title:  nextMusic.Title,
					Artist: nextMusic.Artist,
					Album:  nextMusic.Album,
				},
			}
			musicFound.Occorruces = append(musicFound.Occorruces, musicOccorruce)
		}
		musicsFound = append(musicsFound, musicFound)
	}

	return &musicsFound
}
