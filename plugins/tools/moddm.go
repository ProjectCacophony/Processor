package tools

import (
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/eventlog"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleModDM(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Except(events.NewUserError("common.to-few-params"))
		return
	}

	targetUser, err := event.FindMember(events.WithoutFallbackToSelf())
	if err != nil {
		event.Except(err)
		return
	}

	targetGuild, err := event.State().Guild(event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	message := event.MessageCreate.Content
	for _, field := range []string{event.Prefix(), event.OriginalCommand(), event.Fields()[1]} {
		message = strings.Replace(message, field, "", 1)
	}
	message = strings.Trim(strings.TrimSpace(message), "\"")

	if message == "" {
		event.Except(events.NewUserError("common.to-few-params"))
		return
	}

	fullMessage := event.Translate("tools.help.moddm.dm-prefix", "guild", targetGuild) + "\n" + message

	_, err = event.SendComplexDM(targetUser.ID, discord.MessageCodeToMessage(fullMessage))
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("tools.help.moddm.success", "user", targetUser)

	err = eventlog.CreateItem(event.DB(), event.Publisher(), &eventlog.Item{
		GuildID:     targetGuild.ID,
		ActionType:  eventlog.ActionTypeModDM,
		AuthorID:    event.UserID,
		TargetValue: targetUser.ID,
		TargetType:  eventlog.EntityTypeUser,
		Options: []eventlog.ItemOption{
			{
				Key:      "message code",
				NewValue: message,
				Type:     eventlog.EntityTypeMessageCode,
			},
		},
	})
	event.Except(err)
}
