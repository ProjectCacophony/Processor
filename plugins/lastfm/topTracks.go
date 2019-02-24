// nolint: dupl
package lastfm

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Seklfreak/lastfm-go/lastfm"
	"github.com/bwmarrin/discordgo"
	lastfmclient "gitlab.com/Cacophony/Processor/plugins/lastfm/lastfm-client"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleTopTracks(event *events.Event, lastfmClient *lastfm.Api, offset int) {
	fields := event.Fields()[offset:]

	var makeCollage bool
	period, fields := lastfmclient.GetPeriodFromArgs(fields)
	makeCollage, fields = isCollageRequest(fields)

	// get lastFM username to look up
	username, _ := extractUsername(event, p.db, fields, -1)
	if username == "" {
		event.Respond("lastfm.no-user") // nolint: errcheck
		return
	}

	// lookup user
	userInfo, err := lastfmclient.GetUserinfo(lastfmClient, username)
	if err != nil {
		if strings.Contains(err.Error(), "User not found") {
			event.Respond("lastfm.user-not-found", "username", username) // nolint: errcheck
			return
		}
		event.Except(err)
		return
	}

	// get top artists
	var artists []lastfmclient.ArtistData
	artists, err = lastfmclient.GetTopArtists(lastfmClient, userInfo.Username, 10, period)
	if err != nil {
		event.Except(err)
		return
	}

	// if no artists found, post error and stop
	if len(artists) < 1 {
		event.Respond("lastfm.no-scrobbles", "username", username) // nolint: errcheck
		return
	}

	// get basic embed for user
	embed := getLastfmUserBaseEmbed(userInfo)

	// get top tracks
	var tracks []lastfmclient.TrackData
	tracks, err = lastfmclient.GetTopTracks(lastfmClient, userInfo.Username, 10, period)
	if err != nil {
		event.Except(err)
		return
	}

	// if no tracks found, post error and stop
	if len(tracks) < 1 {
		_, err = event.Respond("lastfm.no-scrobbles", "username", username)
		event.Except(err)
		return
	}

	// set embed title
	embed.Author.Name = "lastfm.tracks.embed.title"

	// create collage if requested
	if makeCollage {
		// initialise variables
		imageURLs := make([]string, 9)
		trackNames := make([]string, 9)
		for i, track := range tracks {
			imageURLs[i] = track.ImageURL
			trackNames[i] = track.Name
			if i >= 8 {
				break
			}
		}

		// create the collage
		collageBytes, err := CollageFromURLs(
			p.httpClient,
			imageURLs,
			trackNames,
			900, 900,
			300, 300,
		)
		if err != nil {
			event.Except(err)
			return
		}

		// add collage image to embed
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "attachment://Cacophony-LastFM-Collage.jpg",
		}
		// send collage to discord and stop
		_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   "Cacophony-LastFM-Collage.jpg",
					Reader: bytes.NewReader(collageBytes),
				},
			},
			Embed: &embed,
		}, "userData", userInfo, "period", period)
		event.Except(err)
		return
	}

	// add tracks to embed
	for i, track := range tracks {
		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " +
			event.Translate("lastfm.track.short", "track", track) +
			"\n"
	}

	// add track image to embed if possible
	if tracks[0].ImageURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: tracks[0].ImageURL,
		}
	}

	// send to discord
	_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
		Embed: &embed,
	}, "userData", userInfo, "period", period)
	event.Except(err)
}
