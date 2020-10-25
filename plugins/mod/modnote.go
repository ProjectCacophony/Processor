package mod

import (
	"gitlab.com/Cacophony/Processor/plugins/eventlog"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleModNote(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Except(events.NewUserError("common.to-few-params"))
		return
	}

	targetUser, err := event.FindMember(events.WithoutFallbackToSelf())
	if err != nil {
		event.Except(err)
		return
	}

	message := event.FieldsVariadic(2)
	if message == "" {
		event.Except(events.NewUserError("common.to-few-params"))
		return
	}

	err = eventlog.CreateItem(event.DB(), event.Publisher(), &eventlog.Item{
		GuildID:     event.GuildID,
		ActionType:  eventlog.ActionTypeModNote,
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

	event.Respond("mod.moddm.success", "targetUser", targetUser, "message", message)
}
