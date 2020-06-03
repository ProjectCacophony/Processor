package admin

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func (p *Plugin) handleDo(event *events.Event) {
	newContent := event.MessageCreate.Content
	for _, field := range append([]string{event.Prefix() + event.OriginalCommand()}, event.Fields()[1:2]...) {
		newContent = strings.Replace(newContent, field, "", 1)
	}
	newContent = strings.TrimSpace(newContent)

	p.logger.Info("running command as root",
		zap.String("user_id", event.UserID),
		zap.String("content", newContent),
	)

	event.Respond("admin.do.progress", "content", newContent)

	// fake the author of event, and update the content
	event.MessageCreate.Content = newContent
	event.SuperUser = true
	err, recoverable := p.publisher.Publish(event.Context(), event)
	if err != nil {
		event.Except(err)
		if !recoverable {
			p.logger.Fatal(
				"received unrecoverable error while publishing \"sudo do\" message",
				zap.Error(err),
			)
		}
		return
	}

	event.React("ok")
}
