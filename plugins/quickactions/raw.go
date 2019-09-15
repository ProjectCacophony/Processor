package quickactions

import (
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

const rawMessageEmoji = "quickaction_raw"

func (p *Plugin) rawMessage(event *events.Event) {
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

	_, err = event.SendDM(
		event.UserID,
		"quickactions.raw.content",
		"emoji", event.MessageReactionAdd.Emoji.MessageFormat(),
		"message", message,
		"content", discord.MessageCodeFromMessage(message),
	)
	if err != nil {
		event.ExceptSilent(err)
		return
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
