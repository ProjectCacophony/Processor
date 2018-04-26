package lastfm

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/project-d-collab/SqsProcessor/models"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/mdb"
)

func getLastfmUserBaseEmbed(userInfo dhelpers.LastfmUserData) (embed discordgo.MessageEmbed) {
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

func getLastfmGuildBaseEmbed(guild *discordgo.Guild, listeners int) (embed discordgo.MessageEmbed) {
	// set embed author
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name: guild.Name,
	}
	if guild.Icon != "" {
		embed.Author.IconURL = discordgo.EndpointGuildIcon(guild.ID, guild.Icon)
	}
	// set embed footer
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    "powered by last.fm",
		IconURL: "https://i.imgur.com/p8wijg4.png",
	}
	if listeners > 0 {
		embed.Footer.Text += " | Total Listeners " + humanize.Comma(int64(listeners))
	}
	// set embed colour
	embed.Color = dhelpers.GetDiscordColorFromHex("#d51007")

	return embed
}

func isCollageRequest(args []string) (collage bool, newArgs []string) {
	return dhelpers.SliceContainsLowerExclude(args, []string{"collage", "image"})
}

func isServerRequest(args []string) (collage bool, newArgs []string) {
	return dhelpers.SliceContainsLowerExclude(args, []string{"server"})
}

func getLastFmUsername(userID string) (username string) {
	var entryBucket models.LastFmEntry
	err := mdb.One(
		models.LastFmTable.DB().Find(bson.M{"userid": userID}),
		&entryBucket,
	)

	if err != nil {
		if !mdb.ErrNotFound(err) {
			logger().Errorln("error requesting user", err.Error())
		}
		return ""
	}

	return entryBucket.LastFmUsername
}
