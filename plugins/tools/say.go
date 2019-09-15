package tools

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleSay(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("tools.say.too-few")
		return
	}

	content := event.MessageCreate.Content

	targetChannel, err := event.State().ChannelFromMention(event.GuildID, event.Fields()[1])
	if err != nil {
		event.Except(err)
		return
	}

	messageCode := strings.Replace(content, event.Prefix()+event.OriginalCommand(), "", 1)
	messageCode = strings.Replace(messageCode, event.Fields()[1], "", 1)
	messageCode = strings.TrimSpace(messageCode)

	message := discord.MessageCodeToMessage(messageCode)

	_, err = event.SendComplex(targetChannel.ID, message)
	if err != nil {
		event.Except(err)
		return
	}

	event.React("ok")
}
