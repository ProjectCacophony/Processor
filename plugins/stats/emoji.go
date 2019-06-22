package stats

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func findEmoji(event *events.Event) (*discordgo.Emoji, string, error) {
	if event.Has(permissions.BotAdmin) {
		for _, fieldA := range event.Fields() {
			for _, fieldB := range event.Fields() {

				emoji, err := event.State().Emoji(fieldA, fieldB)
				if err == nil {
					return emoji, fieldA, nil
				}
			}
		}
	}

	for _, field := range event.Fields() {
		emoji, err := event.State().Emoji(event.GuildID, field)
		if err == nil {
			return emoji, event.GuildID, nil
		}
	}

	return nil, "", errors.New("emoji not found")
}

// TODO: extract emoji ID from emoji

func (p *Plugin) handleEmoji(event *events.Event) {
	emoji, guildID, err := findEmoji(event)
	if err != nil {
		if strings.Contains(err.Error(), "emoji not found") {
			event.Respond("stats.emoji.not-found")
			return
		}
		event.Except(err)
		return
	}

	createdAt, err := discordgo.SnowflakeTimestamp(emoji.ID)
	if err != nil {
		event.Except(err)
		return
	}

	guild, err := event.State().Guild(guildID)
	if err != nil {
		event.Except(err)
		return
	}

	emojiURL := emojiURL(emoji.ID, emoji.Animated)

	_, err = event.RespondComplex(
		&discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "stats.emoji.embed.title",
				Description: "stats.emoji.embed.description",
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "stats.emoji.embed.thumbnail.url",
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "stats.emoji.embed.field.created-at.name",
						Value: "stats.emoji.embed.field.created-at.value",
					},
					{
						Name:  "stats.emoji.embed.field.roles.name",
						Value: "stats.emoji.embed.field.roles.value",
					},
				},
			},
		},
		"emoji", emoji,
		"createdAt", createdAt,
		"guild", guild,
		"emojiURL", emojiURL,
	)
	event.Except(err)
}

func emojiURL(emojiID string, animated bool) string {
	if animated {
		return discordgo.EndpointCDN + "emojis/" + emojiID + ".gif"
	}

	return discordgo.EndpointCDN + "emojis/" + emojiID + ".png"
}
