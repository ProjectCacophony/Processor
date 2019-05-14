package quickactions

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

const (
	starMessageEmoji = "disturd" // TODO: better emoji
)

func (p *Plugin) handleReaction(event *events.Event) bool {
	if event.DM() {
		return false
	}

	channel, err := event.State().Channel(event.ChannelID)
	if err != nil {
		return false
	}

	if channel.Type != discordgo.ChannelTypeGuildText {
		return false
	}

	// ignore in DMs
	if event.MessageReactionAdd.Emoji.Name == starMessageEmoji {
		p.starMessage(event)

		return true
	}

	return false
}

func (p *Plugin) starMessage(event *events.Event) {
	// TODO: use some cache to avoid calling API
	message, err := event.Discord().Client.ChannelMessage(
		event.MessageReactionAdd.ChannelID,
		event.MessageReactionAdd.MessageID,
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}
	message.GuildID = event.MessageReactionAdd.GuildID // :(

	// we need to get the channelID to pin the message,
	// so we can just get it once and use it with SendComplex instead of SendComplexDM,
	// to avoid the overhead of looking it up twice
	channelID, err := discord.DMChannel(
		event.Redis(),
		event.Discord(),
		event.MessageReactionAdd.UserID,
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	// copy message into DMs
	// TODO: better message
	messageIDs, err := event.SendComplex(
		channelID,
		&discordgo.MessageSend{
			Content: "quickactions.star.message.content",
		},
		"message",
		message,
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}
	if len(messageIDs) <= 0 {
		return
	}

	// pin new message in DMs
	err = event.Discord().Client.ChannelMessagePin(
		channelID,
		messageIDs[0].ID,
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	// remove reaction if possible
	if permissions.DiscordManageMessages.Match(
		event.State(),
		event.BotUserID,
		event.MessageReactionAdd.ChannelID,
		false,
	) {
		err = event.Discord().Client.MessageReactionRemove(
			event.MessageReactionAdd.ChannelID,
			event.MessageReactionAdd.MessageID,
			event.MessageReactionAdd.Emoji.APIName(),
			event.MessageReactionAdd.UserID,
		)
		if err != nil {
			event.ExceptSilent(err)
			return
		}
	}
}
