package controllers

import (
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

		playlistMusics := []models.Music{}
		for _, track := range playlistFetched.Tracks.Items {
			music := models.Music{}
			database.DB.Where(&models.Music{SpotifyID: track.Track.ID}).FirstOrCreate(&music)
			hasChanged := music.FillFetchedMusic(&track.Track)
			if hasChanged {	database.DB.Save(&music) }

			playlistMusics = append(playlistMusics, music)
		}

		playlist := models.Playlist{}
		database.DB.Where(&models.Playlist{SpotifyID: playlistFetched.ID}).FirstOrCreate(&playlist)
		hasChanged := playlist.FillFetchedPlaylist(playlistFetched, &playlistMusics)
		if hasChanged {	database.DB.Save(&playlist) }

		for seq, music := range playlist.Musics {
			trackInPlaylist := models.PlaylistMusics{
				PlaylistID: int(playlist.ID),
				MusicID:    int(music.ID),
			}
			database.DB.First(&trackInPlaylist)
			if trackInPlaylist.ReprodutionSequence != seq + 1 {
				trackInPlaylist.SetMusicSequenceInPlaylist(seq + 1)
				database.DB.Save(&trackInPlaylist)
			}

		}

		DatabasePlaylists = append(DatabasePlaylists, &playlist)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"message":   "Playlists atualizadas no banco local",
		"playlists": DatabasePlaylists,
	})

}
