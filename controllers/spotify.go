package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	spotify "github.com/natantn/SpotPlayMe/integrations/spotify"
	"github.com/natantn/SpotPlayMe/services"
)

func GetToken(c *gin.Context) {
	token := spotify.GenerateNewToken()
	c.JSON(200, token)
}

func SyncPlaylists(c *gin.Context) {
	spotifyContext := services.GetSpotifyContext()
	DatabasePlaylists, err := services.FetchPlaylists(spotifyContext)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"message":   "Playlists atualizadas no banco local",
		"playlists": DatabasePlaylists,
	})

}
