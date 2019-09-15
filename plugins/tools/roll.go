package tools

import (
	"math/rand"
	"strconv"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleRoll(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("tools.roll.too-few")
		return
	}

	number, err := strconv.Atoi(event.Fields()[1])
	if err != nil {
		event.Respond("tools.roll.too-few")
		return
	}
	if number < 1 {
		event.Respond("tools.roll.not-positive")
		return
	}

	pick := rand.Intn(number) + 1

	_, err = event.Respond("tools.roll.result", "pick", pick)
	event.Except(err)
}
