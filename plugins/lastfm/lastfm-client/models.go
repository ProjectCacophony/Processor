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

// // LastfmArtistData contains information about an Artist on LastFM
// type LastfmArtistData struct {
// 	Name      string
// 	URL       string
// 	ImageURL  string
// 	Scrobbles int
// }
//
// // LastfmAlbumData contains information about an Album on LastFM
// type LastfmAlbumData struct {
// 	Name      string
// 	URL       string
// 	ImageURL  string
// 	Artist    string
// 	ArtistURL string
// 	Scrobbles int
// }
//
// // LastFmGuildTopTracks contains the top tracks for a guild, it is built by the Worker and stored in redis
// type LastFmGuildTopTracks struct {
// 	GuildID       string
// 	NumberOfUsers int
// 	Period        LastFmPeriod
// 	Tracks        []TrackData
// 	CachedAt      time.Time
// }
