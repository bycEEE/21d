package main

import (
	"net/url"
	"time"
)

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
	SongContributors PrivateSongContributors `json:"SNG_CONTRIBUTORS"`
	LyricsID int `json:"LYRICS_ID"`
	PhysicalReleaseDate string `json:"PHYSICAL_RELEASE_DATE"`
	Copywrite string `json:"COPYWRITE"`
	//BPM string `json:"BPM"`
	Gain string `json:"GAIN"`
	ReleaseDate string `json:"release_date"`
}

// PrivateSongContributors contain additional artist information.
type PrivateSongContributors struct {
	ComposerLyricist []string `json:"composerlyricist"`
	FeaturedArtist []string `json:"featuredartist"`
	MainArtist []string `json:"mainartist"`
	Mixer []string `json:"mixer"`
	Producer []string `json:"producer"`
	StudioPersonnel []string `json:"studiopersonnel"`
	Composer []string `json:"composer"`
	MusicPublisher []string `json:"musicpublisher"`
	Engineer []string `json:"engineer"`
	Writer []string `json:"writer"`
	Author []string `json:"author"`
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

type PublicResults struct {
	Data []PublicTrack `json:"data"`
}

type PublicTrack struct {
	ID int `json:"id"`
	Readable bool `json:"readable"`
	Title string `json:"title"`
	TitleShort string `json:"title_short"`
	TitleVersion string `json:"title_version"`
	Unseen bool `json:"unseen"`
	ISRC string `json:"isrc"`
	Link url.URL `json:"url"`
	Share url.URL `json:"share"`
	Duration int `json:"duration"`
	TrackPosition int `json:"track_position"`
	DiskNumber int `json:"disk_number"`
	Rank int `json:"int"`
	ReleaseDate time.Time `json:"release_date"`
	ExplicitLyrics bool `json:"explicit_lyrics"`
	Preview url.URL `json:"preview"`
	BPM float32 `json:"bpm"`
	Gain float32 `json:"gain"`
	AvailableCountries []string `json:"available_countries"`
	//Alternative PublicTrack `json:"alternative"`
	Contributors []PublicContributor `json:"contributors"`
	Artist PublicArtist `json:"artist"`
	Album PublicAlbum `json:"album"`
	Type string `json:"type,omitempty"`
	Role string `json:"role,omitempty"`
}

type PublicArtist struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Link string `json:"url"`
	Share url.URL `json:"share"`
	Picture url.URL `json:"picture"`
	PictureSmall url.URL `json:"picture_small"`
	PictureMedium url.URL `json:"picture_medium"`
	PictureBig url.URL `json:"picture_big"`
	PictureXL url.URL `json:"picture_xl"`
	NbAlbum int `json:"nb_album"`
	NbFan int `json:"nb_fan"`
	Radio bool `json:"radio"`
	TrackList url.URL `json:"tracklist"`
	Type string `json:"type,omitempty"`
	Role string `json:"role,omitempty"`
}

type PublicAlbum struct {
	ID int `json:"id"`
	Title string `json:"title"`
	UPC string `json:"UPC"`
	Link string `json:"url"`
	Share url.URL `json:"share"`
	Cover url.URL `json:"cover"`
	CoverSmall url.URL `json:"cover_small"`
	CoverMedium url.URL `json:"cover_medium"`
	CoverBig url.URL `json:"cover_big"`
	CoverXL url.URL `json:"cover_xl"`
	GenreID int `json:"genre_id"`
	Genres PublicGenre `json:"genres"`
	Label string `json:"label"`
	NbTracks int `json:"nb_tracks"`
	Duration int `json:"duration"`
	Fans int `json:"fans"`
	Rating int `json:"rating"`
	ReleaseDate time.Time `json:"release_date"`
	RecordType string `json:"record_type"`
	Available bool `json:"available"`
	//Alternative PublicAlbum `json:"alternative"`
	TrackList url.URL `json:"tracklist"`
	ExplicitLyrics bool `json:"explicit_lyrics"`
	Contributors []PublicContributor `json:"contributors"`
	Artist PublicArtist `json:"artist"`
	Type string `json:"type,omitempty"`
	Role string `json:"role,omitempty"`
}

type PublicContributor struct {
	PublicArtist
	Type string `json:"type,omitempty"`
	Role string `json:"role,omitempty"`
}

type PublicGenre struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Picture url.URL `json:"picture"`
	PictureSmall url.URL `json:"picture_small"`
	PictureMedium url.URL `json:"picture_medium"`
	PictureBig url.URL `json:"picture_big"`
	PictureXL url.URL `json:"picture_xl"`
}