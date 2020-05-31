package admin

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
	"go.uber.org/zap"
)

func (p *Plugin) handleAs(event *events.Event) {
	asUser, err := p.state.UserFromMention(event.Fields()[2])
	if err != nil {
		event.Except(err)
		return
	}

	// can not impersonate yourself
	if asUser.ID == event.UserID {
		event.Respond("admin.as.not-yourself")
		return
	}
	// can not impersonate bot admins
	if permissions.BotAdmin.Match(event.State(), event.DB(), asUser.ID, event.ChannelID, event.DM(), false) {
		event.Respond("admin.as.not-as-botadmin")
		return
	}

	newContent := event.MessageCreate.Content
	for _, field := range append([]string{event.Prefix() + event.OriginalCommand()}, event.Fields()[1:3]...) {
		newContent = strings.Replace(newContent, field, "", 1)
	}
	newContent = strings.TrimSpace(newContent)

	p.logger.Info("impersonating user",
		zap.String("user_id", event.UserID),
		zap.String("as_user_id", asUser.ID),
		zap.String("content", newContent),
	)

	event.Respond("admin.as.progress", "content", newContent, "user", asUser)

	// fake the author of event, and update the content
	event.UserID = asUser.ID
	event.MessageCreate.Author = asUser
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
