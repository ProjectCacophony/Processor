package stats

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func findServer(event *events.Event) (*discordgo.Guild, error) {
	if event.Has(permissions.BotAdmin) {
		for _, field := range event.Fields() {
			server, err := event.State().Guild(field)
			if err == nil {
				return server, nil
			}
		}
	}

	return event.State().Guild(event.GuildID)
}

type channelsCount struct {
	Text     int
	Voice    int
	Category int
	Other    int
	Total    int
}

func (p *Plugin) handleServer(event *events.Event) {
	server, err := findServer(event)
	if err != nil {
		if err == redis.Nil {
			event.Respond("stats.server.not-found")
			return
		}
		event.Except(err)
		return
	}

	owner, err := event.State().User(server.OwnerID)
	if err != nil {
		event.Except(err)
		return
	}

	createdAt, err := discordgo.SnowflakeTimestamp(server.ID)
	if err != nil {
		event.Except(err)
		return
	}

	var emojiAnimated int
	for _, emoji := range server.Emojis {
		if emoji.Animated {
			emojiAnimated++
		}
	}

	var channelsCount channelsCount
	for _, channel := range server.Channels {
		switch channel.Type {
		case discordgo.ChannelTypeGuildText:
			channelsCount.Text++
		case discordgo.ChannelTypeGuildVoice:
			channelsCount.Voice++
		case discordgo.ChannelTypeGuildCategory:
			channelsCount.Category++
		default:
			channelsCount.Other++
		}

		channelsCount.Total++
	}

	_, err = event.RespondComplex(&discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "stats.server.embed.title",
			Description: "stats.server.embed.description",
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "stats.server.embed.thumbnail.url",
			},
			Image: &discordgo.MessageEmbedImage{
				URL: "stats.server.embed.image.url",
			},
			URL: "stats.server.embed.url",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "stats.server.embed.field.created-at.name",
					Value: "stats.server.embed.field.created-at.value",
				},
				{
					Name:  "stats.server.embed.field.owner.name",
					Value: "stats.server.embed.field.owner.value",
				},
				{
					Name:  "stats.server.embed.field.members.name",
					Value: "stats.server.embed.field.members.value",
				},
				{
					Name:  "stats.server.embed.field.roles.name",
					Value: "stats.server.embed.field.roles.value",
				},
				{
					Name:  "stats.server.embed.field.emoji.name",
					Value: "stats.server.embed.field.emoji.value",
				},
				{
					Name:  "stats.server.embed.field.channels.name",
					Value: "stats.server.embed.field.channels.value",
				},
				{
					Name:  "stats.server.embed.field.features.name",
					Value: "stats.server.embed.field.features.value",
				},
				{
					Name:  "stats.server.embed.field.nitro-boost.name",
					Value: "stats.server.embed.field.nitro-boost.value",
				},
			},
		},
	},
		"server", server,
		"owner", owner,
		"createdAt", createdAt,
		"iconURL", server.IconURL()+"?size=2048",
		"splashURL", discordgo.EndpointGuildSplash(server.ID, server.Splash)+"?size=2048",
		"bannerURL", discordgo.EndpointGuildBanner(server.ID, server.Banner)+"?size=2048",
		"emojiAnimated", emojiAnimated,
		"channelsCount", channelsCount,
	)
	event.Except(err)
}
