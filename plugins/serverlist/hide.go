package serverlist

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleHide(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("serverlist.hide.too-few-args")
		return
	}

	server := extractExistingServerFromArg(p.redis, p.db, event.Discord(), event.Fields()[2])
	if server == nil {
		event.Respond("serverlist.hide.no-server")
		return
	}

	if !stringSliceContains(event.UserID, server.EditorUserIDs) {
		event.Respond("serverlist.hide.no-editor")
		return
	}

	if server.State != StatePublic && server.State != StateHidden {
		event.Respond("serverlist.hide.wrong-status")
		return
	}

	var err error
	var message string
	if server.State == StateHidden {
		err = server.Unhide(p)
		message = "serverlist.hide.unhide-success"
	} else {
		err = server.Hide(p)
		message = "serverlist.hide.success"
	}
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond(message, "entry", server)
	event.Except(err)
}
