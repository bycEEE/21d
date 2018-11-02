package main

// PrivateResponse contains the response body to a private API call.
type PrivateResponse struct {
	//Error   PrivateError     `json:"error"`
	Results PrivateResults `json:"results"`
}

// PrivateResults contains the results from the private API.
type PrivateResults struct {
	User           PrivateUser `json:"USER"`
	CheckFormLogin string      `json:"checkFormLogin"`
	CheckForm      string      `json:"checkForm"`
	SessionID   string `json:"SESSION_ID"`
	UserToken   string `json:"USER_TOKEN"`
	Country     string `json:"COUNTRY"`
	PlayerToken string `json:"PLAYER_TOKEN"`
	Data PrivateData `json:"DATA"`
	Lyrics PrivateLyrics `json:"LYRICS"`
	//ISRC PrivateISRC `json:"ISRC"`
}

// PrivateTrack is a custom resource with both PrivateResults.Data and PrivateResults.Lyrics
type PrivateTrack struct {
	Data PrivateData `json:"DATA"`
	Lyrics PrivateLyrics `json:"LYRICS"`
}

// PrivateUser contains information about the Deezer user retrieved from the private API.
type PrivateUser struct {
	ID          int    `json:"USER_ID"`
	UserPicture string `json:"USER_PICTURE"`
}

// PrivateData contains information about the resource retrieved from the private API.
type PrivateData struct {
	SongID string `json:"SNG_ID"`
	UploadID int `json:"UPLOAD_ID"`
	SongTitle string `json:"SNG_TITLE"`
	ArtistID string `json:"ART_ID"`
	Artists []PrivateArtist
	AlbumID string `json:"ALB_ID"`
	AlbumTitle string `json:"ALB_TITLE"`
	Type int `json:"TYPE"`
	MD5Origin string `json:"MD5_ORIGIN"`
	AlbumPicture string `json:"ALB_PICTURE"`
	ArtistPicture string `json:"ART_PICTURE"`
	SongRank string `json:"RANK_SNG"`
	// this one is an int for some reason
	FileSizeAAC64 int64 `json:"FILESIZE_AAC_64"`
	FileSizeMP364 string `json:"FILESIZE_MP3_64"`
	FileSizeMP3128 string `json:"FILESIZE_MP3_128"`
	// this one is an int for some reason
	FileSizeMP3256 int64 `json:"FILESIZE_MP3_256"`
	FileSizeMP3320 string `json:"FILESIZE_MP3_320"`
	FileSizeFLAC string `json:"FILESIZE_FLAC"`
	MediaVersion string `json:"MEDIA_VERSION"`
	DiskNumber string `json:"DISK_NUMBER"`
	TrackNumber string `json:"TRACK_NUMBER"`
	Version string `json:"VERSION"`
	ExplicitLyrics string `json:"EXPLICIT_LYRICS"`
	ISRC string `json:"ISRC"`
	//SongContributors PrivateSongContributors SNG_CONTRIBUTORS
	LyricsID int `json:"LYRICS_ID"`
	PhysicalReleaseDate string `json:"PHYSICAL_RELEASE_DATE"`
}

// PrivateLyrics contains the lyrics in text and timestamped form.
type PrivateLyrics struct {
	ID string `json:"LYRICS_ID"`
	LyricsSync []PrivateLyricsSync `json:"LYRICS_SYNC_JSON"`
	LyricsText string `json:"LYRICS_TEXT"`
}

// PrivateLyricsSync contains the lyrics in timestamped form.
type PrivateLyricsSync struct {
	LRCTimestamp string `json:"lrc_timestamp"`
	Milliseconds string `json:"milliseconds"`
	Duration string `json:"duration"`
	Line string `json:"line"`
}

type PrivateArtist struct {
	ID string `json:"ART_ID"`
	RoleID string `json:"ROLE_ID"`
	Name string `json:"ART_NAME"`
	Picture string `json:"ART_PICTURE"`
	Rank string `json:"RANK"`
}

// PrivateError is a mapping of error messages.
type PrivateError struct {
	GatewayError string `json:"GATEWAY_ERROR"`
	ValidTokenRequired string `json:"VALID_TOKEN_REQUIRED"`
	RequestError string `JSON:"REQUEST_ERROR"`
}
