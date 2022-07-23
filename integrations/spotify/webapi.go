package integrations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/natantn/SpotPlayMe/dtos"
)

var (
	SpotifyContext *Spotify
)

type Spotify struct {
	Username         string
	AccessToken      *AccessToken
	PlaylistsFetched []*dtos.PlaylistApiResponse
}

type AccessToken struct {
	Token     string    `json:"access_token"`
	Type      string    `json:"token_type"`
	GeneratedAt time.Time `json:"generated_at"`
	ExpiresIn int       `json:"expires_in"`
}

func GenerateNewToken() *AccessToken {
	clientId := os.Getenv("spotify_client_id")
	clientSecret := os.Getenv("spotify_client_secret")

	client := &http.Client{}
	path := "https://accounts.spotify.com/api/token"
	body := url.Values{}
	body.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(body.Encode()))
	if err != nil {
		panic(err.Error())
	}
	req.SetBasicAuth(clientId, clientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	token := &AccessToken{}
	json.NewDecoder(resp.Body).Decode(token)
	token.GeneratedAt = time.Now()

	return token
}

func (s *Spotify) FetchPlaylistsFromUser() {
	var playlists []dtos.PlaylistItemResponse
	var playlistsApiResponse dtos.UserPlaylistsApiReponse

	webApiRequest := func(path string) *http.Response {
		if path == "" {
			path = fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", s.Username)
		}
		bearerToken := fmt.Sprintf("Bearer %s", s.AccessToken.Token)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			panic(err.Error())
		}
		req.Header.Set("Authorization", bearerToken)

		res, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}

		return res
	}

	path := ""
	for {
		if playlistsApiResponse.Href != "" && playlistsApiResponse.Next != "" {
			path = playlistsApiResponse.Next
			playlistsApiResponse = dtos.UserPlaylistsApiReponse{}
		}

		apiResponse := webApiRequest(path)
		json.NewDecoder(apiResponse.Body).Decode(&playlistsApiResponse)

		playlists = append(playlists, playlistsApiResponse.Items...)

		if playlistsApiResponse.Next == "" {
			break
		}
	}

	for _, p := range playlists {
		s.GetPlaylistById(p.ID)
	}
}

func (s *Spotify) GetPlaylistById(id string) {

	playlist := dtos.PlaylistApiResponse{}
	webPlaylistApiRequest := func(id string) *http.Response {
		path := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s", id)
		bearerToken := fmt.Sprintf("Bearer %s", s.AccessToken.Token)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, path, nil)
		if err != nil {
			panic(err.Error())
		}
		req.Header.Set("Authorization", bearerToken)

		res, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}

		return res
	}

	apiResponse := webPlaylistApiRequest(id)
	json.NewDecoder(apiResponse.Body).Decode(&playlist)

	s.PlaylistsFetched = append(s.PlaylistsFetched, &playlist)
}
