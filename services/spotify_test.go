package services

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	database "github.com/natantn/SpotPlayMe/integrations/database"
	"github.com/stretchr/testify/assert"
)

func setupEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err.Error())
	}

	database.ConectDB()
}

func TestSpotifyContextNonExisting(t *testing.T) {
	setupEnv()
	spotifyContext := GetSpotifyContext()

	assert.NotNil(t, spotifyContext, "Spotify context é nulo")
}

func TestSpotifyContextExpiredToken(t *testing.T) {
	setupEnv()

	spotifyContext := GetSpotifyContext()
	firstToken := spotifyContext.AccessToken
	spotifyContext.AccessToken.GeneratedAt = spotifyContext.AccessToken.GeneratedAt.Add(time.Duration(-1) * time.Hour)

	spotifyContext = GetSpotifyContext()
	secondToken := spotifyContext.AccessToken

	assert.NotEqual(t, secondToken.Token, firstToken.Token, "Token são iguais")
}
