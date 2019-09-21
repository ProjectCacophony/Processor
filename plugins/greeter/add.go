package greeter

import (
	"fmt"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleAdd(event *events.Event, greeterType greeterType) {
	if len(event.Fields()) < 3 {
		event.Respond("common.to-few-params")
		return
	}

	targetChannel, err := event.FindChannel()
	if err != nil {
		event.Except(err)
		return
	}

	var message string
	if len(event.Fields()) >= 4 {
		message = event.Fields()[3]
	}

	existingEntry, err := entryFind(event.DB(), event.GuildID, targetChannel.ID, greeterType)
	if err != nil {
		event.Except(err)
		return
	}

	if existingEntry == nil ||
		existingEntry.Rule.ID == 0 ||
		len(existingEntry.Rule.Actions) < 1 ||
		len(existingEntry.Rule.Actions[0].Values) < 1 {
		if message == "" {
			event.Respond("common.to-few-params")
			return
		}

		rule := &models.Rule{
			GuildID: event.GuildID,
			Name:    fmt.Sprintf("Greeter: %s #%s", greeterType, targetChannel.Name),
			Actions: []models.RuleAction{
				{
					Name: "send_message_to",
					Values: []string{
						message,
						targetChannel.ID,
					},
				},
			},
			Silent:  true,
			Managed: true,
		}

		switch greeterType {
		case greeterTypeJoin:
			rule.TriggerName = "when_join"
		case greeterTypeLeave:
			rule.TriggerName = "when_leave"
		default:
			event.Except(fmt.Errorf("unknown greeter type: %d", greeterType))
			return
		}

		err = models.CreateRule(event.DB(), rule)
		if err != nil {
			event.Except(err)
			return
		}

		err = entryAdd(
			event.DB(),
			event.GuildID,
			targetChannel.ID,
			greeterType,
			message,
			rule,
		)
		if err != nil {
			// clean up created rule
			models.DeleteRule(event.DB(), rule)
			event.Except(err)
			return
		}
	} else {
		if message == "" {
			err = models.DeleteRule(event.DB(), &existingEntry.Rule)
			if err != nil {
				event.Except(err)
				return
			}

			err = entryDelete(event.DB(), existingEntry.ID)
			if err != nil {
				event.Except(err)
				return
			}

			event.Respond("greeter.add.remove-success", "greeterType", greeterType, "channel", targetChannel)
			return
		}

		err = entryUpdate(event.DB(), existingEntry.ID, message)
		if err != nil {
			event.Except(err)
			return
		}

		existingEntry.Rule.Actions[0].Values[0] = message
		err = models.UpdateRule(event.DB(), existingEntry.RuleID, &existingEntry.Rule)
		if err != nil {
			event.Except(err)
			return
		}
	}

	event.Respond("greeter.add.success", "greeterType", greeterType, "channel", targetChannel)
}
