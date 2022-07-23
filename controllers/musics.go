package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/natantn/SpotPlayMe/services"
)

func SearchMusic(c *gin.Context) {

	title := strings.ReplaceAll(c.Query("title"), "%20", " ")
	title = "%" + title + "%"

	musics := services.GetMusicsByTitle(title)
	if len(*musics) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": "Música não encontrada em playlists cadastradas no banco",
		})
		return
	}

	musicsFound := services.GetMusicOccorrencesInPlaylists(musics)
	c.JSON(http.StatusOK, musicsFound)
	return

}
