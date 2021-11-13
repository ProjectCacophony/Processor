package tools

import (
	"strings"
	"time"

	// embed timezones
	_ "time/tzdata"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleTime(event *events.Event) {
	kst, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("tools.time.content",
		"now", time.Now(),
		"utc", time.FixedZone("UTC", 0),
		"kst", kst,
	)
}

func (p *Plugin) handleTimezone(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("tools.timezone.too-few")
		return
	}

	zone, err := time.LoadLocation(event.Fields()[1])
	if err != nil {
		if strings.Contains(err.Error(), "unknown time zone") {
			event.Respond("tools.timezone.not-found")
			return
		}

		event.Except(err)
		return
	}

	err = event.SetTimezone(zone)
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("tools.timezone.done")
}
