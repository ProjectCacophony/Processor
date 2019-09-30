package tools

import (
	"math/rand"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleChoose(event *events.Event) {
	parts := event.Fields()[1:]
	items := make([]string, 0, len(parts))
	for _, field := range parts {
		if field == "|" {
			continue
		}

		items = append(items, field)
	}

	if len(items) <= 0 {
		event.Respond("tools.choose.too-few")
		return
	}

	pick := items[rand.Intn(len(items))]

	_, err := event.Respond("tools.choose.result", "pick", pick)
	event.Except(err)
}
