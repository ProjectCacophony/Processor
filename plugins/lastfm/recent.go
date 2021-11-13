package lastfm

import (
	"strings"

	"github.com/Seklfreak/lastfm-go/lastfm"
	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
	lastfmclient "gitlab.com/Cacophony/Processor/plugins/lastfm/lastfm-client"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleRecent(event *events.Event, lastfmClient *lastfm.Api) {
	fields := event.Fields()[2:]

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

	// get recent tracks
	var tracks []lastfmclient.TrackData
	tracks, err = lastfmclient.GetRecentTracks(lastfmClient, userInfo.Username, 100)
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
	embed.Author.Name = "lastfm.recent.embed.title"

	var embeds []*discordgo.MessageEmbed

	// add tracks to embed
	for i, track := range tracks {
		embed.Description += event.Translate("lastfm.track.long", "track", track, "hidenp", true)

		if track.NowPlaying {
			embed.Description += " - _" + event.Translate("lastfm.recent.nowplaying") + "_"
		} else if !track.Time.IsZero() {
			embed.Description += " - " + humanize.Time(track.Time)
		}

		embed.Description += "\n"

		if i > 0 && (i+1)%10 == 0 {
			send := discord.TranslateMessageSend(
				event.Localizations(),
				&discordgo.MessageSend{
					Embed: embed,
				},
				"userData", userInfo,
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
