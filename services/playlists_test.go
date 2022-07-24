package services

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/natantn/SpotPlayMe/dtos"
	integrations "github.com/natantn/SpotPlayMe/integrations/database"
	"github.com/stretchr/testify/assert"
)

func mockedItemsToFetch() (*dtos.TrackItemInPlaylistResponse, *dtos.PlaylistApiResponse) {
	trackItemMock := dtos.TrackItemInPlaylistResponse{
		Name: "Mocked track",
		Artists: []struct {
			Name string "json:\"name\""
		}{
			{Name: "The Mocker"},
		},
		Album: struct {
			Name        string "json:\"name\""
			ReleaseDate string "json:\"release_date\""
		}{
			Name: "Mocking Mocked Mocks",
		},
	}
	playlistMock := dtos.PlaylistApiResponse{
		Name:        "Mocked Playlist",
		Description: "Testing if this mocked playlist is going to syncronized in database",
		Tracks: struct {
			Items []struct {
				Track dtos.TrackItemInPlaylistResponse "json:\"track\""
			} "json:\"items\""
		}{
			Items: []struct {
				Track dtos.TrackItemInPlaylistResponse "json:\"track\""
			}{
				{trackItemMock},
			},
		},
	}

	return &trackItemMock, &playlistMock
}

func TestSycNoPlaylists(t *testing.T) {
	setupEnv()
	spotifyContext := GetSpotifyContext()

	spotifyContext.PlaylistsFetched = []*dtos.PlaylistApiResponse{}

	playlists, err := FetchPlaylists(spotifyContext)

	if assert.Error(t, err, "Deveria retornar um erro") {
		expectedMessage := "Não foram encontradas playlists para sincronização no usuário fornecido"
		assert.ErrorContains(t, err, expectedMessage, fmt.Sprintf("O erro deveria ser '%s'", expectedMessage))
	}
	assert.Nil(t, playlists, "Não deveria haver playlists sincronizadas")

}

func TestSyncNewPlaylist(t *testing.T) {
	setupEnv()
	spotifyContext := GetSpotifyContext()

	trackMocked, playlistMocked := mockedItemsToFetch()

	spotifyContext.PlaylistsFetched = []*dtos.PlaylistApiResponse{playlistMocked}

	playlists, err := FetchPlaylists(spotifyContext)

	assert.NoError(t, err)

	playlistResult := playlists[0]
	musicsResult := playlistResult.Musics

	assert.NotEmpty(t, playlistResult)
	assert.Equal(t, playlistMocked.Name, playlistResult.Title)
	assert.NotEmpty(t, musicsResult)
	assert.Equal(t, trackMocked.Name, musicsResult[0].Title)

	integrations.DB.Delete(&playlistResult)
}

func TestSyncExistingPlaylist(t *testing.T) {
	setupEnv()
	spotifyContext := GetSpotifyContext()

	trackMocked, playlistMocked := mockedItemsToFetch()

	spotifyContext.PlaylistsFetched = []*dtos.PlaylistApiResponse{playlistMocked}
	playlists, err := FetchPlaylists(spotifyContext)

	trackMockedUpdated := dtos.TrackItemInPlaylistResponse{
		Name: "Mocked track - Demo",
		Artists: trackMocked.Artists,
		Album: trackMocked.Album,
	}
	playlistMockedUpdated := *playlistMocked
	playlistMockedUpdated.Name = "The ultimate mocked playlist"
	playlistMockedUpdated.ID = strconv.FormatUint(uint64(playlists[0].ID), 10)
	playlistMockedUpdated.Tracks.Items = []struct{Track dtos.TrackItemInPlaylistResponse "json:\"track\""}{{trackMockedUpdated}}
	spotifyContext.PlaylistsFetched = []*dtos.PlaylistApiResponse{&playlistMockedUpdated}

	playlistsUpdated, err := FetchPlaylists(spotifyContext)

	assert.NoError(t, err)

	playlistResult := playlistsUpdated[0]
	musicsResult := playlistResult.Musics

	assert.NotEmpty(t, playlistResult)
	assert.NotEqual(t, playlistMocked.Name, playlistResult.Title)
	assert.Equal(t, playlistMockedUpdated.Name, playlistResult.Title)
	assert.NotEmpty(t, musicsResult)
	assert.NotEqual(t, trackMocked.Name, musicsResult[0].Title)
	assert.Equal(t, trackMockedUpdated.Name, musicsResult[0].Title)

	integrations.DB.Delete(&playlistResult)
}
