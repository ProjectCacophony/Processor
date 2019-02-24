package lastfm

import (
	"strings"

	"github.com/Seklfreak/lastfm-go/lastfm"
	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
	lastfmclient "gitlab.com/Cacophony/Processor/plugins/lastfm/lastfm-client"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleRecent(event *events.Event, lastfmClient *lastfm.Api) {
	fields := event.Fields()[2:]

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

	// get recent tracks
	var tracks []lastfmclient.TrackData
	tracks, err = lastfmclient.GetRecentTracks(lastfmClient, userInfo.Username, 10)
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

	// add tracks to embed
	for _, track := range tracks {
		embed.Description += event.Translate("lastfm.track.long", "track", track, "hidenp", true)

		if track.NowPlaying {
			embed.Description += " - _" + event.Translate("lastfm.recent.nowplaying") + "_"
		} else if !track.Time.IsZero() {
			embed.Description += " - " + humanize.Time(track.Time)
		}

		embed.Description += "\n"
	}

	// send to discord
	_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
		Embed: &embed,
	}, "userData", userInfo)
	event.Except(err)
}
