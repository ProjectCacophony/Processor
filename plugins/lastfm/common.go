package lastfm

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	lastfm_client "gitlab.com/Cacophony/Processor/plugins/lastfm/lastfm-client"
	"gitlab.com/Cacophony/go-kit/discord"
)

// getLastfmUserBaseEmbed gets a discordgo embed base for a last.fm user
func getLastfmUserBaseEmbed(userInfo lastfm_client.UserData) (embed discordgo.MessageEmbed) {
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
		embed.Footer.Text += " | Total Plays " + strconv.Itoa(userInfo.Scrobbles) // TODO: humanize
	}
	// set embed colour
	embed.Color = discord.HexToColorCode("#d51007")

	return embed
}

// // getLastfmGuildBaseEmbed gets a discordgo embed base for a discord guild
// func getLastfmGuildBaseEmbed(guild *discordgo.Guild, listeners int) (embed discordgo.MessageEmbed) {
// 	// set embed author
// 	embed.Author = &discordgo.MessageEmbedAuthor{
// 		Name: guild.Name,
// 	}
// 	if guild.Icon != "" {
// 		embed.Author.IconURL = discordgo.EndpointGuildIcon(guild.ID, guild.Icon)
// 	}
// 	// set embed footer
// 	embed.Footer = &discordgo.MessageEmbedFooter{
// 		Text:    "powered by last.fm",
// 		IconURL: "https://i.imgur.com/p8wijg4.png",
// 	}
// 	if listeners > 0 {
// 		embed.Footer.Text += " | Total Listeners " + strconv.Itoa(listeners) // TODO: humanize
// 	}
// 	// set embed colour
// 	embed.Color = discord.HexToColorCode("#d51007")
//
// 	return embed
// }

// TODO: add username persistency
// // getLastFmUsername gets the lastFM username for a specific discord userID, returns an empty string if none found
// func getLastFmUsername(ctx context.Context, userID string) (username string) {
// 	var entryBucket models.LastFmEntry
// 	err := models.LastFmRepository.FindOne(
// 		ctx,
// 		bson.NewDocument(bson.EC.String("userid", userID)),
// 		&entryBucket,
// 	)
//
// 	if err != nil {
// 		if err == mongo.ErrNotFound {
// 			return ""
// 		}
// 		logger().Errorln("error requesting user", err.Error())
// 	}
//
// 	return entryBucket.LastFmUsername
// }
