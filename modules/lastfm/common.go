package lastfm

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
	"gitlab.com/project-d-collab/dhelpers"
)

func getLastfmBaseEmbed(userInfo dhelpers.LastfmUserData) (embed discordgo.MessageEmbed) {
	// set embed author
	embed.Author = &discordgo.MessageEmbedAuthor{
		URL: "https://www.last.fm/user/" + userInfo.Username,
	}
	if userInfo.Icon != "" {
		embed.Author.IconURL = userInfo.Icon
	}
	// set embed footer
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    "powered by last.fm",
		IconURL: "https://i.imgur.com/p8wijg4.png",
	}
	if userInfo.Scrobbles > 0 {
		embed.Footer.Text += " | Total Plays " + humanize.Comma(int64(userInfo.Scrobbles))
	}
	// set embed colour
	embed.Color = dhelpers.GetDiscordColorFromHex("#d51007")

	return embed
}

func isCollageRequest(args []string) (collage bool) {
	for _, arg := range args {
		arg = strings.ToLower(arg)
		if arg == "collage" || arg == "image" {
			return true
		}
	}
	return false
}
