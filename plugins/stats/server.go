package stats

import (
	"github.com/bwmarrin/discordgo"
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

func (p *Plugin) handleServer(event *events.Event) {
	server, err := findServer(event)
	if err != nil {
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
		"iconURL", discordgo.EndpointGuildIcon(server.ID, server.Icon)+"?size=1024",
		"splashURL", discordgo.EndpointGuildSplash(server.ID, server.Splash)+"?size=1024",
		"bannerURL", discordgo.EndpointGuildBanner(server.ID, server.Banner)+"?size=1024",
	)
	event.Except(err)
}
