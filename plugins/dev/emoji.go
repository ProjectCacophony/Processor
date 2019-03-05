package dev

import (
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/discord/emoji"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleDevEmoji(event *events.Event) {
	// translating and sending message manually, to avoid automatic replacing of the emoji codes
	// should usually be avoided!
	message := event.Translate("dev.emoji.list", "emojiList", emoji.List)

	pages := discord.Pagify(message)

	for _, page := range pages {
		_, err := event.Discord().Client.ChannelMessageSend(event.MessageCreate.ChannelID, page)
		if err != nil {
			event.Except(err)
			return
		}
	}
}
