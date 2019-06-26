package quickactions

import (
	"gitlab.com/Cacophony/go-kit/events"
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
	)
	if err != nil {
		event.ExceptSilent(err)
		return
	}

	// TODO: add support for embeds?
	// TODO: replace emoji (do not display IDs)?
}
