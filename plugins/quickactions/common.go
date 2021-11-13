package quickactions

import (
	"github.com/bwmarrin/discordgo"
)

// TODO: put into go-kit
func convertMessageToEmbed(message *discordgo.Message) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{}

	if message != nil {
		embed.Description = message.Content

		if len(message.Embeds) > 0 {
			embed.Description += "\n" + message.Embeds[0].Description

			embed.URL = message.Embeds[0].URL
			embed.Timestamp = message.Embeds[0].Timestamp
			embed.Color = message.Embeds[0].Color
			embed.Footer = message.Embeds[0].Footer
			embed.Image = message.Embeds[0].Image
			embed.Thumbnail = message.Embeds[0].Thumbnail
			embed.Video = message.Embeds[0].Video
			embed.Provider = message.Embeds[0].Provider
			embed.Author = message.Embeds[0].Author
			embed.Fields = message.Embeds[0].Fields
		}
	}

	return embed
}
