package quickactions

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleReaction(event *events.Event) bool {
	if event.DM() {
		return false
	}

	channel, err := event.State().Channel(event.ChannelID)
	if err != nil {
		return false
	}

	// following actions work in DMs too

	if _, ok := remindEmoji[event.MessageReactionAdd.Emoji.Name]; ok {
		p.remindMessage(event)

		return true
	}

	if event.MessageReactionAdd.Emoji.Name == rawMessageEmoji {
		p.rawMessage(event)

		return true
	}

	// ignore in DMs
	if channel.Type != discordgo.ChannelTypeGuildText {
		return false
	}

	// following actions do NOT work in DMs

	if event.MessageReactionAdd.Emoji.Name == starMessageEmoji {
		p.starMessage(event)

		return true
	}

	return false
}
