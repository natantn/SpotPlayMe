package dtos

// Spotify API
type PlaylistApiResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Tracks      struct {
		Items []struct {
			Track TrackItemInPlaylistResponse `json:"track"`
		} `json:"items"`
	} `json:"tracks"`
}

type TrackItemInPlaylistResponse struct {
	Album struct {
		Name        string `json:"name"`
		ReleaseDate string `json:"release_date"`
	} `json:"album"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PlaylistItemResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          string `json:"id"`
}

type UserPlaylistsApiReponse struct {
	Href     string                 `json:"href"`
	Items    []PlaylistItemResponse `json:"items"`
	Limit    int                    `json:"limit"`
	Next     string                 `json:"next"`
	Offset   int                    `json:"offset"`
	Previous string                 `json:"previous"`
}
