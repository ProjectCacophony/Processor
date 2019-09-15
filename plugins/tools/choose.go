package tools

import (
	"math/rand"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleChoose(event *events.Event) {
	var items []string // nolint: prealloc
	for _, field := range event.Fields()[1:] {
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
