package template

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func handleFoo(event *events.Event) {
	_, err := event.Respond("changeme.foo.response")
	if err != nil {
		event.Except(err)
	}
}

func handleBar(event *events.Event) {
	_, err := event.Respond("changeme.bar.response")
	if err != nil {
		event.Except(err)
	}
}
