package lastfmclient

import (
	"strconv"
	"time"

	"github.com/Seklfreak/lastfm-go/lastfm"
)

// GetUserinfo returns information about a LastFM user
func GetUserinfo(client *lastfm.Api, username string) (userData UserData, err error) {
	// request data
	var lastfmUser lastfm.UserGetInfo
	lastfmUser, err = client.User.GetInfo(lastfm.P{"user": username})
	if err != nil {
		return userData, err
	}
	// parse fields into lastfmUserData
	userData.Username = lastfmUser.Name
	userData.Name = lastfmUser.RealName
	userData.Country = lastfmUser.Country
	if lastfmUser.PlayCount != "" {
		userData.Scrobbles, _ = strconv.Atoi(lastfmUser.PlayCount) // nolint: errcheck, gas
	}

	if len(lastfmUser.Images) > 0 {
		for _, image := range lastfmUser.Images {
			if image.Size == lastFmTargetImageSize {
				userData.Icon = image.Url
			}
		}
	}

	if lastfmUser.Registered.Unixtime != "" {
		timeI, err := strconv.ParseInt(lastfmUser.Registered.Unixtime, 10, 64)
		if err == nil {
			userData.AccountCreation = time.Unix(timeI, 0)
		}
	}

	return userData, nil
}

// GetRecentTracks returns recent tracks listened to by an user
func GetRecentTracks(client *lastfm.Api, username string, limit int) (tracksData []TrackData, err error) {
	// request data
	var lastfmRecentTracks lastfm.UserGetRecentTracksExtended
	lastfmRecentTracks, err = client.User.GetRecentTracksExtended(lastfm.P{
		"limit": limit + 1, // in case nowplaying + already scrobbled
		"user":  username,
	})
	if err != nil {
		return nil, err
	}

	// parse fields
	if lastfmRecentTracks.Total > 0 {
		for i, track := range lastfmRecentTracks.Tracks {
			if i == 1 {
				// prevent nowplaying + already scrobbled
				if lastfmRecentTracks.Tracks[0].Url == track.Url {
					continue
				}
			}
			lastTrack := TrackData{
				Name:      track.Name,
				URL:       track.Url,
				Artist:    track.Artist.Name,
				ArtistURL: track.Artist.Url,
				Album:     track.Album.Name,
				Loved:     false,
			}
			for _, image := range track.Images {
				if image.Size == lastFmTargetImageSize {
					lastTrack.ImageURL = image.Url
				}
			}
			for _, image := range track.Artist.Image {
				if image.Size == lastFmTargetImageSize {
					lastTrack.ArtistImageURL = image.Url
				}
			}
			if track.Loved == "1" || track.Loved == "true" {
				lastTrack.Loved = true
			}
			if track.NowPlaying == "1" || track.NowPlaying == "true" {
				lastTrack.NowPlaying = true
			}

			timestamp, err := strconv.Atoi(track.Date.Uts)
			if err == nil {
				lastTrack.Time = time.Unix(int64(timestamp), 0)
			}

			tracksData = append(tracksData, lastTrack)
			if len(tracksData) >= limit {
				break
			}
		}
	}

	return tracksData, nil
}

// nolint
// // LastFmGetTopArtists returns the top artists of an user
// func LastFmGetTopArtists(lastfmClient *lastfm.Api, lastfmUsername string, limit int, period LastFmPeriod) (artistsData []LastfmArtistData, err error) {
// 	// request data
// 	var lastfmTopArtists lastfm.UserGetTopArtists
// 	lastfmTopArtists, err = lastfmClient.User.GetTopArtists(lastfm.P{
// 		"limit":  limit,
// 		"user":   lastfmUsername,
// 		"period": string(period),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// parse fields
// 	if lastfmTopArtists.Total > 0 {
// 		for _, artist := range lastfmTopArtists.Artists {
// 			lastArtist := LastfmArtistData{
// 				Name: artist.Name,
// 				URL:  artist.Url,
// 			}
// 			for _, image := range artist.Images {
// 				if image.Size == lastFmTargetImageSize {
// 					lastArtist.ImageURL = image.Url
// 				}
// 			}
// 			lastArtist.Scrobbles, _ = strconv.Atoi(artist.PlayCount) // nolint: gas
//
// 			artistsData = append(artistsData, lastArtist)
// 			if len(artistsData) >= limit {
// 				break
// 			}
// 		}
// 	}
//
// 	return artistsData, nil
// }
//
// // LastFmGetTopTracks returns the top tracks of an user
// func LastFmGetTopTracks(lastfmClient *lastfm.Api, lastfmUsername string, limit int, period LastFmPeriod) (tracksData []TrackData, err error) {
// 	// request data
// 	var lastfmTopTracks lastfm.UserGetTopTracks
// 	lastfmTopTracks, err = lastfmClient.User.GetTopTracks(lastfm.P{
// 		"limit":  limit,
// 		"user":   lastfmUsername,
// 		"period": string(period),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// parse fields
// 	if lastfmTopTracks.Total > 0 {
// 		for _, track := range lastfmTopTracks.Tracks {
// 			lastTrack := TrackData{
// 				Name:      track.Name,
// 				URL:       track.Url,
// 				Artist:    track.Artist.Name,
// 				ArtistURL: track.Artist.Url,
// 			}
// 			for _, image := range track.Images {
// 				if image.Size == lastFmTargetImageSize {
// 					lastTrack.ImageURL = image.Url
// 				}
// 			}
// 			lastTrack.Scrobbles, _ = strconv.Atoi(track.PlayCount) // nolint: gas
//
// 			tracksData = append(tracksData, lastTrack)
// 			if len(tracksData) >= limit {
// 				break
// 			}
// 		}
// 	}
//
// 	return tracksData, nil
// }
//
// // LastFmGetTopAlbums returns the top albums of an user
// func LastFmGetTopAlbums(lastfmClient *lastfm.Api, lastfmUsername string, limit int, period LastFmPeriod) (albumsData []LastfmAlbumData, err error) {
// 	// request data
// 	var lastfmTopAlbums lastfm.UserGetTopAlbums
// 	lastfmTopAlbums, err = lastfmClient.User.GetTopAlbums(lastfm.P{
// 		"limit":  limit,
// 		"user":   lastfmUsername,
// 		"period": string(period),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// parse fields
// 	if lastfmTopAlbums.Total > 0 {
// 		for _, album := range lastfmTopAlbums.Albums {
// 			lastAlbum := LastfmAlbumData{
// 				Name:      album.Name,
// 				URL:       album.Url,
// 				Artist:    album.Artist.Name,
// 				ArtistURL: album.Artist.Url,
// 			}
// 			for _, image := range album.Images {
// 				if image.Size == lastFmTargetImageSize {
// 					lastAlbum.ImageURL = image.Url
// 				}
// 			}
// 			lastAlbum.Scrobbles, _ = strconv.Atoi(album.PlayCount) // nolint: gas
//
// 			albumsData = append(albumsData, lastAlbum)
// 			if len(albumsData) >= limit {
// 				break
// 			}
// 		}
// 	}
//
// 	return albumsData, nil
// }
//
// // LastFmGetPeriodFromArgs parses args to figure out the correct period
// func LastFmGetPeriodFromArgs(args []string) (LastFmPeriod, []string) {
// 	for i, arg := range args {
// 		switch arg {
// 		case "7day", "7days", "week", "7", "seven":
// 			args = append(args[:i], args[i+1:]...)
// 			return LastFmPeriod7day, args
// 		case "1month", "month", "1", "one":
// 			args = append(args[:i], args[i+1:]...)
// 			return LastFmPeriod1month, args
// 		case "3month", "threemonths", "quarter", "3", "three":
// 			args = append(args[:i], args[i+1:]...)
// 			return LastFmPeriod3month, args
// 		case "6month", "halfyear", "half", "sixmonths", "6", "six":
// 			args = append(args[:i], args[i+1:]...)
// 			return LastFmPeriod6month, args
// 		case "12month", "year", "twelvemonths", "12", "twelve":
// 			args = append(args[:i], args[i+1:]...)
// 			return LastFmPeriod12month, args
// 		case "overall", "all", "alltime", "all-time":
// 			args = append(args[:i], args[i+1:]...)
// 			return LastFmPeriodOverall, args
// 		}
// 	}
//
// 	return LastFmPeriodOverall, args
// }
