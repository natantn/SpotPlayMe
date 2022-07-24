package services

import (
	"os"
	"time"

	spotify "github.com/natantn/SpotPlayMe/integrations/spotify"
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
