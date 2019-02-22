package dev

import (
	"strconv"
	"time"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleDevSleep(event *events.Event) {
	var secondsToSleep int
	if len(event.Fields()) >= 3 {
		secondsToSleep, _ = strconv.Atoi(event.Fields()[2])
	}

	if secondsToSleep <= 0 || secondsToSleep > 60 {
		secondsToSleep = 10
	}

	durationToSleep := time.Duration(secondsToSleep) * time.Second

	_, err := event.Respond("dev.sleep.start", "durationToSleep", durationToSleep)
	if err != nil {
		event.Except(err)
		return
	}

	time.Sleep(durationToSleep)

	_, err = event.Respond("dev.sleep.finished")
	if err != nil {
		event.Except(err)
		return
	}
}
