package dev

import (
	"errors"
	"fmt"
	"time"

	"gitlab.com/Cacophony/go-kit/events"
)

const questionnaireKey = "cacophony:processor:dev:questionnaire"

func (p *Plugin) handleDevQuestionnaire(event *events.Event) {
	err := event.Questionnaire().Register(
		questionnaireKey,
		events.QuestionnaireFilter{
			GuildID:   event.GuildID,
			ChannelID: event.ChannelID,
			UserID:    event.UserID,
			Type:      events.MessageCreateType,
		},
		map[string]interface{}{
			"started": time.Now(),
		},
	)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("I will be waiting for your response")
	event.Except(err)
}

func (p *Plugin) handleDevQuestionnaireMatch(event *events.Event) {
	startedAtData, ok := event.QuestionnaireMatch.Payload["started"].(string)
	if !ok {
		event.Except(errors.New("received invalid questionnaire match payload"))
		return
	}

	var startedAt time.Time
	err := startedAt.UnmarshalText([]byte(startedAtData))
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Send(
		event.ChannelID,
		fmt.Sprintf(
			"<@%s> You took %s to respond!",
			event.UserID,
			time.Since(startedAt).Round(time.Millisecond)),
	)
	event.Except(err)
}
