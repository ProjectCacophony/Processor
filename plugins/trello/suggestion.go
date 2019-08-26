package trello

import (
	"strings"

	trello "github.com/VojtechVitek/go-trello"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleSuggestion(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("common.invalid-params")
		return
	}

	issueTitle := event.Fields()[1]

	// if there are more than 3 fields, combine all of them to make the title to avoid require quotes
	if len(event.Fields()) > 3 {
		issueTitle = strings.Join(event.Fields()[1:], " ")
	}

	if len(issueTitle) > 50 {
		event.Respond("trello.title.to-long")
		return
	}

	issueDescription := ""
	if len(event.Fields()) == 3 {
		issueDescription = event.Fields()[2]
	}

	list, err := p.trello.List(backlogBoardID)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = list.AddCard(trello.Card{
		Name: strings.TrimSpace(issueTitle),
		Desc: strings.TrimSpace(issueDescription),
	})
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("trello.suggestion.received")
}
