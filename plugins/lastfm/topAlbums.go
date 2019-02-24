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

func (p *Plugin) handleTopAlbums(event *events.Event, lastfmClient *lastfm.Api, offset int) {
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

	// get basic embed for user
	embed := getLastfmUserBaseEmbed(userInfo)

	// get top albums
	var albums []lastfmclient.AlbumData
	albums, err = lastfmclient.GetTopAlbums(lastfmClient, userInfo.Username, 10, period)
	if err != nil {
		event.Except(err)
		return
	}

	// if no albums found, post error and stop
	if len(albums) < 1 {
		_, err = event.Respond("lastfm.no-scrobbles", "username", username)
		event.Except(err)
		return
	}

	// set content
	embed.Author.Name = "lastfm.albums.embed.title"

	// create collage if requested
	if makeCollage {
		// initialise variables
		imageURLs := make([]string, 9)
		albumNames := make([]string, 9)
		for i, album := range albums {
			imageURLs[i] = album.ImageURL
			albumNames[i] = album.Name
			if i >= 8 {
				break
			}
		}

		// create the collage
		collageBytes, err := CollageFromURLs(
			p.httpClient,
			imageURLs,
			albumNames,
			900, 900,
			300, 300,
		)
		if err != nil {
			event.Except(err)
			return
		}

		// add collage image to embed
		embed.Image = &discordgo.MessageEmbedImage{
			URL: "attachment://Cacophony-LastFM-Collage.png",
		}
		// send collage to discord and stop
		_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				{
					Name:   "Cacophony-LastFM-Collage.png",
					Reader: bytes.NewReader(collageBytes),
				},
			},
			Embed: &embed,
		}, "userData", userInfo, "period", period)
		event.Except(err)
		return
	}

	// add albums to embed
	for i, album := range albums {
		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " +
			event.Translate("lastfm.album.short", "album", album) +
			"\n"
	}

	// add album image to embed if possible
	if albums[0].ImageURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: albums[0].ImageURL,
		}
	}

	// send to discord
	_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
		Embed: &embed,
	}, "userData", userInfo, "period", period)
	event.Except(err)
}
