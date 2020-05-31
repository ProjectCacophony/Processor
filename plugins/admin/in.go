package admin

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func (p *Plugin) handleIn(event *events.Event) {
	inChannel, err := p.state.ChannelFromMentionTypesEverywhere(event.Fields()[2], discordgo.ChannelTypeGuildText)
	if err != nil {
		event.Except(err)
		return
	}

	// can not run in current channel
	if inChannel.ID == event.ChannelID {
		event.Respond("admin.in.not-current")
		return
	}

	newContent := event.MessageCreate.Content
	for _, field := range append([]string{event.Prefix() + event.OriginalCommand()}, event.Fields()[1:3]...) {
		newContent = strings.Replace(newContent, field, "", 1)
	}
	newContent = strings.TrimSpace(newContent)

	p.logger.Info("running command in channel remotely",
		zap.String("user_id", event.UserID),
		zap.String("in_channel_id", inChannel.ID),
		zap.String("content", newContent),
	)

	event.Send(
		inChannel.ID,
		"admin.in.progress",
		"content", newContent,
		"user", event.MessageCreate.Author,
	)

	// fake the author of event, and update the content
	event.MessageCreate.ChannelID = inChannel.ID
	event.MessageCreate.Content = newContent
	err, recoverable := p.publisher.Publish(event.Context(), event)
	if err != nil {
		event.Except(err)
		if !recoverable {
			p.logger.Fatal(
				"received unrecoverable error while publishing \"sudo as\" message",
				zap.Error(err),
			)
		}
		return
	}

	event.React("ok")
}
