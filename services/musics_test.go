package services

import (
	"testing"

	integrations "github.com/natantn/SpotPlayMe/integrations/database"
	"github.com/natantn/SpotPlayMe/models"
	"github.com/stretchr/testify/assert"
)

func TestGetNotExistingMusic(t *testing.T) {
	setupEnv()
	musics := GetMusicsByTitle("Lorem Impsum")
	assert.Equal(t, 0, len(*musics), "Música não existente foi encontrada")
}

func TestGetExistingMusics(t *testing.T) {
	setupEnv()
	music := models.Music{
		Title: "Mocked Music",
	}
	integrations.DB.Create(&music)

	musicFound := GetMusicsByTitle(music.Title)
	assert.Equal(t, 1, len(*musicFound))

	integrations.DB.Delete(&music)
}

func TestGetOccurrenceOfMusic(t *testing.T) {
	setupEnv()
	musics := []models.Music{
		{
			Title:  "Music A",
			Artist: "Mocked",
		},
		{
			Title:  "Music X",
			Artist: "Mocked",
		},
	}
	playlist := models.Playlist{
		Title:  "Mock",
		Musics: musics,
	}
	integrations.DB.Create(&playlist)

	occurrences := GetMusicOccorrencesInPlaylists(&musics)

	assert.Equal(t, 2, len(*occurrences))

	integrations.DB.Delete(&playlist)
}

func TestGetOccurrenceOfMusicWithNoOccurrence(t *testing.T) {
	setupEnv()
	musics := []models.Music{
		{
			Title:  "Music A",
			Artist: "Mocked",
		},
		{
			Title:  "Music X",
			Artist: "Mocked",
		},
	}
	playlist := models.Playlist{
		Title:  "Mock",
		Musics: musics,
	}
	integrations.DB.Create(&playlist)

	occorruence := models.PlaylistMusics{
		PlaylistID: int(playlist.ID),
		MusicID:    int(musics[0].ID),
	}
	integrations.DB.Delete(&occorruence)

	occurrences := GetMusicOccorrencesInPlaylists(&musics)

	assert.Equal(t, 2, len(*occurrences))

	integrations.DB.Delete(&playlist)
}
