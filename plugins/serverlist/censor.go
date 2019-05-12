package serverlist

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleCensor(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("serverlist.censor.too-few-args")
		return
	}

	server := extractExistingServerFromArg(p.redis, p.db, event.Discord(), event.Fields()[2])
	if server == nil {
		event.Respond("serverlist.censor.no-server")
		return
	}

	var err error
	var message string
	if server.State == StateCensored {
		err = server.Uncensor(p)
		message = "serverlist.censor.uncensor-success"
	} else {
		if len(event.Fields()) < 4 {
			event.Respond("serverlist.censor.too-few-args")
			return
		}

		err = server.Censor(p, event.Fields()[3])
		message = "serverlist.censor.success"
	}
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond(message, "entry", server)
	event.Except(err)
}
