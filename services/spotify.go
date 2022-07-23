package services

import (
	"errors"
	"os"
	"time"

	database "github.com/natantn/SpotPlayMe/integrations/database"
	spotify "github.com/natantn/SpotPlayMe/integrations/spotify"
	"github.com/natantn/SpotPlayMe/models"
)

func GetSpotifyContext() *spotify.Spotify {
	ctx := spotify.SpotifyContext

	if ctx == nil {
		username := os.Getenv("spotify_username")
		spotify.SpotifyContext = &spotify.Spotify{
			Username:    username,
			AccessToken: spotify.GenerateNewToken(),
		}
		ctx = spotify.SpotifyContext
	}

	if time.Now().After(ctx.AccessToken.GeneratedAt.Add(time.Second * time.Duration(ctx.AccessToken.ExpiresIn))) {
		ctx.AccessToken = spotify.GenerateNewToken()
	}

	return ctx
}

func FetchPlaylists(s *spotify.Spotify) (DatabasePlaylists []*models.Playlist, err error) {

	s.FetchPlaylistsFromUser()
	if len(s.PlaylistsFetched) == 0 {
		DatabasePlaylists = nil
		err = errors.New("Não foram encontradas playlists para sincronização no usuário fornecido")
	}

	for _, playlistFetched := range s.PlaylistsFetched {

		//Fetching each music in playlist
		playlistMusics := []models.Music{}
		for _, track := range playlistFetched.Tracks.Items {
			music := models.Music{}
			database.DB.Where(&models.Music{SpotifyID: track.Track.ID}).FirstOrCreate(&music)
			hasChanged := music.FillFetchedMusic(&track.Track)
			if hasChanged {
				database.DB.Save(&music)
			}
			playlistMusics = append(playlistMusics, music)
		}

		//Fetching playlist
		playlist := models.Playlist{}
		database.DB.Where(&models.Playlist{SpotifyID: playlistFetched.ID}).FirstOrCreate(&playlist)
		hasChanged := playlist.FillFetchedPlaylist(playlistFetched, &playlistMusics)
		if hasChanged {
			database.DB.Save(&playlist)
		}

		//Updating musics sequence
		for seq, music := range playlist.Musics {
			trackInPlaylist := models.PlaylistMusics{
				PlaylistID: int(playlist.ID),
				MusicID:    int(music.ID),
			}
			database.DB.First(&trackInPlaylist)
			if trackInPlaylist.ReprodutionSequence != seq+1 {
				trackInPlaylist.SetMusicSequenceInPlaylist(seq + 1)
				database.DB.Save(&trackInPlaylist)
			}
		}
		DatabasePlaylists = append(DatabasePlaylists, &playlist)
	}

	return
}
