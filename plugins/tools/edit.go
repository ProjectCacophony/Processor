package tools

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleEdit(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("tools.say.too-few")
		return
	}

	targetMessage, err := event.FindMessageLink(event.Fields()[1])
	if err != nil {
		event.Except(err)
		return
	}

	messageCode := strings.Replace(event.MessageCreate.Content, event.Prefix()+event.OriginalCommand(), "", 1)
	messageCode = strings.Replace(messageCode, event.Fields()[1], "", 1)
	messageCode = strings.TrimSpace(messageCode)

	message := discord.MessageCodeToMessage(messageCode)

	edit := &discordgo.MessageEdit{
		Content: &message.Content,
		Embed:   message.Embed,
		ID:      targetMessage.ID,
		Channel: targetMessage.ChannelID,
	}

	_, err = discord.EditComplexWithVars(event.Redis(), event.Discord(), event.Localizations(), edit, event.DM())
	if err != nil {
		event.Except(err)
		return
	}

	event.React("ok")
}
