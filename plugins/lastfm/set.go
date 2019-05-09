package lastfm

import (
	lastfmclient "gitlab.com/Cacophony/Processor/plugins/lastfm/lastfm-client"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleSet(event *events.Event) {
	if len(event.Fields()) < 3 {
		return
	}

	username := event.Fields()[2]

	_, err := lastfmclient.GetUserinfo(p.lastfmClient, username)
	if err != nil {
		event.Respond("lastfm.user-not-found", "username", username)
		return
	}

	// upsert username to db
	err = setLastFmUsername(p.db, event.MessageCreate.Author.ID, username)
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
