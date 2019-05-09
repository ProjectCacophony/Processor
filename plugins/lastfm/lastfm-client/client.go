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
		userData.Scrobbles, _ = strconv.ParseInt(lastfmUser.PlayCount, 10, 64)
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

// GetTopArtists returns the top artists of an user
func GetTopArtists(client *lastfm.Api, username string, limit int, period Period) ([]ArtistData, error) {
	// request data
	var artistsData []ArtistData
	var err error
	var lastfmTopArtists lastfm.UserGetTopArtists
	lastfmTopArtists, err = client.User.GetTopArtists(lastfm.P{
		"limit":  limit,
		"user":   username,
		"period": string(period),
	})
	if err != nil {
		return nil, err
	}

	// parse fields
	if lastfmTopArtists.Total > 0 {
		for _, artist := range lastfmTopArtists.Artists {
			lastArtist := ArtistData{
				Name: artist.Name,
				URL:  artist.Url,
			}
			for _, image := range artist.Images {
				if image.Size == lastFmTargetImageSize {
					lastArtist.ImageURL = image.Url
				}
			}
			lastArtist.Scrobbles, _ = strconv.Atoi(artist.PlayCount)

			artistsData = append(artistsData, lastArtist)
			if len(artistsData) >= limit {
				break
			}
		}
	}

	return artistsData, nil
}

// GetTopTracks returns the top tracks of an user
func GetTopTracks(client *lastfm.Api, username string, limit int, period Period) ([]TrackData, error) {
	var tracksData []TrackData
	var err error
	// request data
	var lastfmTopTracks lastfm.UserGetTopTracks
	lastfmTopTracks, err = client.User.GetTopTracks(lastfm.P{
		"limit":  limit,
		"user":   username,
		"period": string(period),
	})
	if err != nil {
		return nil, err
	}

	// parse fields
	if lastfmTopTracks.Total > 0 {
		for _, track := range lastfmTopTracks.Tracks {
			lastTrack := TrackData{
				Name:      track.Name,
				URL:       track.Url,
				Artist:    track.Artist.Name,
				ArtistURL: track.Artist.Url,
			}
			for _, image := range track.Images {
				if image.Size == lastFmTargetImageSize {
					lastTrack.ImageURL = image.Url
				}
			}
			lastTrack.Scrobbles, _ = strconv.Atoi(track.PlayCount)

			tracksData = append(tracksData, lastTrack)
			if len(tracksData) >= limit {
				break
			}
		}
	}

	return tracksData, nil
}

// LastFmGetTopAlbums returns the top albums of an user
func GetTopAlbums(client *lastfm.Api, username string, limit int, period Period) ([]AlbumData, error) {
	var albumsData []AlbumData
	var err error
	// request data
	var lastfmTopAlbums lastfm.UserGetTopAlbums
	lastfmTopAlbums, err = client.User.GetTopAlbums(lastfm.P{
		"limit":  limit,
		"user":   username,
		"period": string(period),
	})
	if err != nil {
		return nil, err
	}

	// parse fields
	if lastfmTopAlbums.Total > 0 {
		for _, album := range lastfmTopAlbums.Albums {
			lastAlbum := AlbumData{
				Name:      album.Name,
				URL:       album.Url,
				Artist:    album.Artist.Name,
				ArtistURL: album.Artist.Url,
			}
			for _, image := range album.Images {
				if image.Size == lastFmTargetImageSize {
					lastAlbum.ImageURL = image.Url
				}
			}
			lastAlbum.Scrobbles, _ = strconv.Atoi(album.PlayCount)

			albumsData = append(albumsData, lastAlbum)
			if len(albumsData) >= limit {
				break
			}
		}
	}

	return albumsData, nil
}

// GetPeriodFromArgs parses args to figure out the correct period
func GetPeriodFromArgs(args []string) (Period, []string) {
	for i, arg := range args {
		switch arg {
		case "7day", "7days", "week", "7", "seven":
			args = append(args[:i], args[i+1:]...)
			return Period7day, args
		case "1month", "month", "1", "one":
			args = append(args[:i], args[i+1:]...)
			return Period1month, args
		case "3month", "threemonths", "quarter", "3", "three":
			args = append(args[:i], args[i+1:]...)
			return Period3month, args
		case "6month", "halfyear", "half", "sixmonths", "6", "six":
			args = append(args[:i], args[i+1:]...)
			return Period6month, args
		case "12month", "year", "twelvemonths", "12", "twelve":
			args = append(args[:i], args[i+1:]...)
			return Period12month, args
		case "overall", "all", "alltime", "all-time":
			args = append(args[:i], args[i+1:]...)
			return PeriodOverall, args
		}
	}

	return PeriodOverall, args
}
