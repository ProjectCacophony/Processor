package lastfm

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func handleSet(event *events.Event) {
	if len(event.Fields()) < 3 {
		return
	}

	username := event.Fields()[2]

	// TODO: first validate if user exists

	// upsert username to db
	err := setLastFmUsername(event.DB(), event.MessageCreate.Author.ID, username)
	if err != nil {
		event.Except(err)
		return
	}

	// send to discord
	_, err = event.Respond("lastfm.set.saved")
	if err != nil {
		event.Except(err)
		return
	}
}
