package integrations

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
}

func (t *AccessToken) GenerateToken() {
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
	
	json.NewDecoder(resp.Body).Decode(t)
}
