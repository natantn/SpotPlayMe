package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	database "github.com/natantn/SpotPlayMe/integrations/database"
	"github.com/natantn/SpotPlayMe/models"
)

func SearchMusic(c *gin.Context) {

	title := c.Query("title")
	title = strings.ReplaceAll(title, "%20", " ")
	title = "%" + title + "%"

	musics := []models.Music{}
	database.DB.Where("title LIKE ?", title).Find(&musics)

	if len(musics) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Música não encontrada em playlists cadastradas no banco",
		})
		return
	}

	type MusicItem struct {
		Title  string
		Artist string
		Album  string
	}

	type MusicOccorruce struct {
		PlaylistTitle        string    `json:"playlist_title"`
		ReproduntionSequence int       `json:"reprodution_sequence"`
		Previous             MusicItem `json:"previous_music"`
		Next                 MusicItem `json:"next_music"`
	}

	type MusicFound struct {
		Music      MusicItem        `json:"music"`
		Occorruces []MusicOccorruce `json:"occurrences"`
	}

	musicsFound := []MusicFound{}

	for _, music := range musics {
		occurrences := []models.PlaylistMusics{}
		database.DB.Where(&models.PlaylistMusics{MusicID: int(music.ID)}).Find(&occurrences)

		musicFound := MusicFound{
			Music: MusicItem{music.Title, music.Artist, music.Album},
		}

		for _, occurrence := range occurrences {
			playlist := models.Playlist{}
			database.DB.First(&playlist, occurrence.PlaylistID)

			seq := occurrence.ReprodutionSequence

			previousMusicReprodution := models.PlaylistMusics{
				PlaylistID:          occurrence.PlaylistID,
				ReprodutionSequence: seq - 1,
			}
			previousMusic := models.Music{}
			if seq-1 > 0 {
				database.DB.Where(&previousMusicReprodution).Take(&previousMusicReprodution)
				if previousMusicReprodution.MusicID != 0 {
					database.DB.First(&previousMusic, previousMusicReprodution.MusicID)
				}
			}

			nextMusicReprodution := models.PlaylistMusics{
				PlaylistID:          occurrence.PlaylistID,
				ReprodutionSequence: seq + 1,
			}
			nextMusic := models.Music{}
			database.DB.Where(&nextMusicReprodution).Take(&nextMusicReprodution)
			if nextMusicReprodution.MusicID != 0 {
				database.DB.First(&nextMusic, nextMusicReprodution.MusicID)
			}

			musicOccorruce := MusicOccorruce{
				PlaylistTitle:        playlist.Title,
				ReproduntionSequence: occurrence.ReprodutionSequence,
				Previous: MusicItem{
					Title:  previousMusic.Title,
					Artist: previousMusic.Artist,
					Album:  previousMusic.Album,
				},
				Next: MusicItem{
					Title:  nextMusic.Title,
					Artist: nextMusic.Artist,
					Album:  nextMusic.Album,
				},
			}

			musicFound.Occorruces = append(musicFound.Occorruces, musicOccorruce)
		}

		musicsFound = append(musicsFound, musicFound)
	}

	c.JSON(http.StatusOK, musicsFound)
	return

}
