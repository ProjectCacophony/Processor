package tools

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func (p *Plugin) handleSay(event *events.Event) {
	if len(event.Fields()) < 3 &&
		!(len(event.Fields()) >= 2 && len(event.MessageCreate.Attachments) > 0) {
		event.Respond("tools.say.too-few")
		return
	}

	targetChannel, err := event.State().ChannelFromMention(event.GuildID, event.Fields()[1])
	if err != nil {
		event.Except(err)
		return
	}

	messageCode := strings.Replace(event.MessageCreate.Content, event.Prefix()+event.OriginalCommand(), "", 1)
	messageCode = strings.Replace(messageCode, event.Fields()[1], "", 1)
	messageCode = strings.TrimSpace(messageCode)

	message := discord.MessageCodeToMessage(messageCode)

	for _, attachment := range event.MessageCreate.Attachments {
		data, err := event.HTTPClient().Get(attachment.URL)
		if err != nil {
			event.Logger().Error(
				"failure downloading attachment for say post",
				zap.Error(err), zap.String("url", attachment.URL),
			)
			continue
		}
		defer data.Body.Close()

		message.Files = append(message.Files, &discordgo.File{
			Name:   attachment.Filename,
			Reader: data.Body,
		})
	}

	_, err = event.SendComplex(targetChannel.ID, message)
	if err != nil {
		event.Except(err)
		return
	}

	event.React("ok")
}
func (p *Plugin) handleGet(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("tools.get.too-few")
		return
	}

	message, err := event.FindMessageLink(event.Fields()[1])
	if err != nil {
		event.Except(err)
		return
	}

	code := discord.MessageCodeFromMessage(message)

	_, err = event.Respond("tools.get.result", "code", code)
	event.Except(err)
}
