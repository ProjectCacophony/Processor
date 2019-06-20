package stats

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func findAnyChannelEverywhere(event *events.Event) (*discordgo.Channel, error) {
	for _, field := range event.Fields() {
		channel, err := event.State().ChannelFromMentionTypesEverywhere(field)
		if err == nil {
			return channel, nil
		}
	}

	return event.State().Channel(event.ChannelID)
}

func (p *Plugin) handleChannel(event *events.Event) {
	channel, err := event.FindAnyChannel()
	if event.Has(permissions.BotAdmin) {
		channel, err = findAnyChannelEverywhere(event)
	}
	if err != nil {
		event.Except(err)
		return
	}

	createdAt, err := discordgo.SnowflakeTimestamp(channel.ID)
	if err != nil {
		event.Except(err)
		return
	}

	parentChannel, _ := event.State().Channel(channel.ParentID)
	guild, _ := event.State().Guild(channel.GuildID)

	_, err = event.RespondComplex(&discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       "stats.channel.embed.title",
			Description: "stats.channel.embed.description",
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "stats.channel.embed.thumbnail.url",
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "stats.channel.embed.field.created-at.name",
					Value: "stats.channel.embed.field.created-at.value",
				},
				{
					Name:  "stats.channel.embed.field.permission-overwrites.name",
					Value: "stats.channel.embed.field.permission-overwrites.value",
				},
			},
		},
	},
		"channel", channel,
		"createdAt", createdAt,
		"iconURL", discordgo.EndpointGroupIcon(channel.ID, channel.Icon)+"?size=2048",
		"parentChannel", parentChannel,
		"guild", guild,
	)
	event.Except(err)
}
