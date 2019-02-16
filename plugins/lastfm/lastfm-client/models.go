package lastfmclient

import "time"

// UserData contains information about an User on LastFM
type UserData struct {
	Username        string
	Name            string
	Icon            string
	Scrobbles       int64
	Country         string
	AccountCreation time.Time
}

// TrackData contains information about a Track on LastFM
type TrackData struct {
	Name           string
	URL            string
	ImageURL       string
	Artist         string
	ArtistURL      string
	ArtistImageURL string
	Album          string
	Time           time.Time
	Loved          bool
	NowPlaying     bool
	Scrobbles      int
	// used for guild stats
	Users int
}

// ArtistData contains information about an Artist on LastFM
type ArtistData struct {
	Name      string
	URL       string
	ImageURL  string
	Scrobbles int
}

// AlbumData contains information about an Album on LastFM
type AlbumData struct {
	Name      string
	URL       string
	ImageURL  string
	Artist    string
	ArtistURL string
	Scrobbles int
}

// // GuildTopTracks contains the top tracks for a guild, it is built by the Worker and stored in redis
// type GuildTopTracks struct {
// 	GuildID       string
// 	NumberOfUsers int
// 	Period        Period
// 	Tracks        []TrackData
// 	CachedAt      time.Time
// }
