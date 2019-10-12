// nolint: dupl
package lastfm

import (
	"fmt"
	"strings"

	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/paginator"

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
	username := extractUsername(event, p.db, fields, -1)
	if username == "" {
		event.Respond("lastfm.no-user")
		return
	}

	// lookup user
	userInfo, err := lastfmclient.GetUserinfo(lastfmClient, username)
	if err != nil {
		if strings.Contains(err.Error(), "User not found") {
			event.Respond("lastfm.user-not-found", "username", username)
			return
		}
		event.Except(err)
		return
	}

	// get basic embed for user
	embed := getLastfmUserBaseEmbed(userInfo)

	// get top albums
	var albums []lastfmclient.AlbumData
	albums, err = lastfmclient.GetTopAlbums(lastfmClient, userInfo.Username, 200, period)
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
		var files []*paginator.File
		var imageURLs, trackNames []string

		for i, album := range albums {

			imageURLs = append(imageURLs, album.ImageURL)
			trackNames = append(trackNames, album.Name)

			if i > 0 && (i+1)%9 == 0 {
				// create the collage
				collageBytes, err := CollageFromURLs(
					event.HTTPClient(),
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
			event.Localizations(),
			&discordgo.MessageSend{
				Embed: embed,
			},
			"userData", userInfo, "period", period,
		)
		err = event.Paginator().ImagePaginator(
			event.BotUserID,
			event.ChannelID,
			event.UserID,
			send.Embed,
			files,
			event.DM(),
		)
		event.Except(err)
		return
	}

	var embeds []*discordgo.MessageEmbed

	// add albums to embed
	for i, album := range albums {
		embed.Description += fmt.Sprintf("`#%2d`", i+1) + " " +
			event.Translate("lastfm.album.short", "album", album) +
			"\n"

		if i > 0 && (i+1)%10 == 0 {
			if albums[i-9].ImageURL != "" {
				embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
					URL: albums[i-9].ImageURL,
				}
			}

			send := discord.TranslateMessageSend(
				event.Localizations(),
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
		event.BotUserID,
		event.ChannelID,
		event.UserID,
		event.DM(),
		embeds...,
	)
	event.Except(err)
}
