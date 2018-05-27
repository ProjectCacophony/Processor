package lastfm

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/Cacophony/SqsProcessor/models"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/mdb"
	"gitlab.com/Cacophony/dhelpers/slice"
)

// getLastfmUserBaseEmbed gets a discordgo embed base for a last.fm user
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
	embed.Color = dhelpers.HexToDecimal("#d51007")

	return embed
}

// getLastfmGuildBaseEmbed gets a discordgo embed base for a discord guild
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
	embed.Color = dhelpers.HexToDecimal("#d51007")

	return embed
}

// isCollageRequest returns true if the args contain an arg making it a collage request
func isCollageRequest(args []string) (collage bool, newArLastFmNoUserPassedgs []string) {
	return slice.ContainsLowerExclude(args, []string{"collage", "image"})
}

// isServerRequest returns true if the args contain an arg making it a server request
func isServerRequest(args []string) (serverRequest bool, newArgs []string) { // nolint: unparam
	return slice.ContainsLowerExclude(args, []string{"server"})
}

// getLastFmUsername gets the lastFM username for a specific discord userID, returns an empty string if none found
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
