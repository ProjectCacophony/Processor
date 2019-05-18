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

	// ignore in DMs
	if channel.Type != discordgo.ChannelTypeGuildText {
		return false
	}

	if event.MessageReactionAdd.Emoji.Name == starMessageEmoji {
		p.starMessage(event)

		return true
	}

	if _, ok := remindEmoji[event.MessageReactionAdd.Emoji.Name]; ok {
		p.remindMessage(event)

		return true
	}

	return false
}
