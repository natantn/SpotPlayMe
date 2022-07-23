package dtos

type MusicItem struct {
	Title  string
	Artist string
	Album  string
}

type MusicOccorruce struct {
	PlaylistTitle        string    `json:"playlist_title"`
	ReproduntionSequence int       `json:"reprodution_sequence"`
	Previous             MusicItem `json:"previous_music"`
	Next                 MusicItem `json:"next_music"`
}

type MusicFound struct {
	Music      MusicItem        `json:"music"`
	Occorruces []MusicOccorruce `json:"occurrences"`
}
