package greeter

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/Processor/plugins/automod/actions"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/discord"
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
	var autoDelete time.Duration
	if len(event.Fields()) >= 5 {
		duration, err := time.ParseDuration(event.Fields()[4])
		if err == nil && duration > 1*time.Second && duration < 24*time.Hour {
			autoDelete = duration
		}
	}

	existingEntry, err := entryFind(event.DB(), event.GuildID, targetChannel.ID, greeterType)
	if err != nil {
		event.Except(err)
		return
	}

	rule, err := newRule(
		event,
		greeterType,
		targetChannel,
		message,
		autoDelete,
	)
	if err != nil {
		event.Except(events.AsUserError(err))
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
			autoDelete,
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

		err = entryUpdate(event.DB(), existingEntry.ID, message, autoDelete)
		if err != nil {
			event.Except(err)
			return
		}

		existingEntry.Rule = *rule
		err = models.UpdateRule(event.DB(), existingEntry.RuleID, &existingEntry.Rule)
		if err != nil {
			event.Except(err)
			return
		}
	}

	sampleMessageCode := actions.ReplaceText(&models.Env{
		State:     event.State(),
		GuildID:   event.GuildID,
		ChannelID: []string{targetChannel.ID},
		UserID:    []string{event.UserID},
	}, message)
	messageSend := discord.MessageCodeToMessage(sampleMessageCode)

	event.Respond(
		"greeter.add.success",
		"greeterType", greeterType, "channel", targetChannel, "autoDelete", autoDelete,
	)
	event.RespondComplex(messageSend)
}

func newRule(
	event *events.Event,
	greeterType greeterType,
	targetChannel *discordgo.Channel,
	message string,
	autoDelete time.Duration,
) (*models.Rule, error) {
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
	case greeterTypeBan:
		rule.TriggerName = "when_ban"
	case greeterTypeUnban:
		rule.TriggerName = "when_unban"
	default:
		return nil, fmt.Errorf("unknown greeter type: %d", greeterType)
	}

	if autoDelete.Seconds() > 0 {
		rule.Actions = append(
			rule.Actions,
			models.RuleAction{
				Name: "wait",
				Values: []string{
					autoDelete.String(),
				},
			},
			models.RuleAction{
				Name: "delete_bot_message",
			},
		)
	}

	return rule, nil
}
