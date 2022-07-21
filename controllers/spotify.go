package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	database "github.com/natantn/SpotPlayMe/integrations/database"
	spotify "github.com/natantn/SpotPlayMe/integrations/spotify"
	"github.com/natantn/SpotPlayMe/models"
)

func GetToken(c *gin.Context) {
	token := spotify.GenerateToken()
	c.JSON(200, token)
}

func SyncPlaylists(c *gin.Context) {
	username := os.Getenv("spotify_username")
	spotContext := spotify.Spotify{
		Username:    username,
		AccessToken: spotify.GenerateToken(),
	}

	spotContext.FetchPlaylistsFromUser()

	DatabasePlaylists := []*models.Playlist{}

	for _, playlistFetched := range spotContext.PlaylistsFetched {
		playlist := models.Playlist{}

		database.DB.Where(&models.Playlist{SpotifyID: playlist.SpotifyID}).FirstOrCreate(&playlist)

		playlist.Title = playlistFetched.Name
		playlist.Description = playlistFetched.Description
		playlist.SpotifyID = playlistFetched.ID

		for _, track := range playlistFetched.Tracks.Items {
			music := models.Music{}
			database.DB.Where(&models.Music{SpotifyID: track.Track.ID}).FirstOrCreate(&music)

			music.Title = track.Track.Name
			music.Album = track.Track.Album.Name
			music.ReleaseDate = track.Track.Album.ReleaseDate
			for i, artist := range track.Track.Artists {
				music.Artist += fmt.Sprintf(" %s", artist.Name)
				if i < len(track.Track.Artists)-1 {
					music.Artist += ","
				}
			}

			database.DB.Save(&music)
			playlist.Musics = append(playlist.Musics, music)
		}

		database.DB.Save(&playlist)
		DatabasePlaylists = append(DatabasePlaylists, &playlist)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"message":   "Playlists atualizadas no banco local",
		"playlists": DatabasePlaylists,
	})

}
