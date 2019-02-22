package lastfm

import (
	"strings"

	lastfmclient "gitlab.com/Cacophony/Processor/plugins/lastfm/lastfm-client"

	"github.com/Seklfreak/lastfm-go/lastfm"
	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleAbout(event *events.Event, lastfmClient *lastfm.Api) {
	fields := event.Fields()[2:]

	// get lastFM username to look up
	username, _ := extractUsername(event, p.db, fields, -1)

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
	embed.Author.Name = "lastfm.about.embed.title"

	// add fields
	// replace scrobbles count in footer with field
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "🎶 Scrobbles",
		Value:  humanize.Comma(userInfo.Scrobbles),
		Inline: false,
	})
	if strings.Contains(embed.Footer.Text, "|") {
		embed.Footer.Text = strings.TrimSpace(strings.SplitN(embed.Footer.Text, "|", 2)[0])
	}

	// add country to embed if possible
	if userInfo.Country != "" && userInfo.Country != "None" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🗺 Country",
			Value:  userInfo.Country,
			Inline: false,
		})
	}

	// add account creation date to embed if possible
	if !userInfo.AccountCreation.IsZero() {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🗓 Account creation",
			Value:  humanize.Time(userInfo.AccountCreation),
			Inline: false,
		})
	}

	// replace author icon with bigger thumbnail for about
	if userInfo.Icon != "" {
		embed.Author.IconURL = ""
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: userInfo.Icon,
		}
	}

	// send to discord
	_, err = event.SendComplex(event.MessageCreate.ChannelID, &discordgo.MessageSend{
		Embed: &embed,
	}, "userData", userInfo)
	event.Except(err)
}
