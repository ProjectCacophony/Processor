package quickactions

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

const starMessageEmoji = "quickaction_star"

func (p *Plugin) starMessage(event *events.Event) {
	message, err := getMessage(
		event.State(),
		event.Discord(),
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
	messageIDs, err := event.SendComplex(
		channelID,
		&discordgo.MessageSend{
			Content: "quickactions.message.content",
			Embed:   convertMessageToEmbed(message),
		},
		"message",
		message,
		"emoji",
		event.MessageReactionAdd.Emoji.MessageFormat(),
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
		if errD, ok := err.(*discordgo.RESTError); !ok ||
			errD.Message == nil ||
			errD.Message.Code != discordgo.ErrCodeMaximumPinsReached {
			event.ExceptSilent(err)
			return
		}
	}

	// remove reaction if possible
	if permissions.DiscordManageMessages.Match(
		event.State(),
		p.db,
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
