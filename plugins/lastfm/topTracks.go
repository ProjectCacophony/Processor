// nolint: dupl
package lastfm

import (
	"fmt"
	"strings"

	"gitlab.com/Cacophony/go-kit/paginator"

	"gitlab.com/Cacophony/go-kit/discord"

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

	// get basic embed for user
	embed := getLastfmUserBaseEmbed(userInfo)

	// get top tracks
	var tracks []lastfmclient.TrackData
	tracks, err = lastfmclient.GetTopTracks(lastfmClient, userInfo.Username, 200, period)
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
		var files []*paginator.File
		var imageURLs, trackNames []string

		for i, track := range tracks {

			imageURLs = append(imageURLs, track.ImageURL)
			trackNames = append(trackNames, track.Name)

			if i > 0 && (i+1)%9 == 0 {
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
				files = append(files, &paginator.File{
					Name: "Cacophony-LastFM-Collage.jpg",
					Data: collageBytes,
				})

				imageURLs = []string{}
				trackNames = []string{}
			}
			if len(files) >= 3 {
				break
			}
		}

		var send = discord.TranslateMessageSend(
			event.Localisations(),
			&discordgo.MessageSend{
				Embed: embed,
			},
			"userData", userInfo, "period", period,
		)
		err = event.Paginator().ImagePaginator(
			event.GuildID,
			event.ChannelID,
			event.UserID,
			send.Embed,
			files,
		)
		event.Except(err)
		return
	}

	var embeds []*discordgo.MessageEmbed

	// add tracks to embed
	for i, track := range tracks {
		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " +
			event.Translate("lastfm.track.short", "track", track) +
			"\n"

		if i > 0 && (i+1)%10 == 0 {
			if tracks[i-9].ImageURL != "" {
				embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
					URL: tracks[i-9].ImageURL,
				}
			}

			send := discord.TranslateMessageSend(
				event.Localisations(),
				&discordgo.MessageSend{
					Embed: embed,
				},
				"userData", userInfo, "period", period,
			)

			tempEmbed := *send.Embed
			embeds = append(embeds, &tempEmbed)

			embed.Description = ""
		}
	}

	err = event.Paginator().EmbedPaginator(
		event.GuildID,
		event.ChannelID,
		event.UserID,
		embeds...,
	)
	event.Except(err)
}
