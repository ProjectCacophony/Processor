package serverlist

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleHide(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("serverlist.hide.too-few-args") // nolint: errcheck
		return
	}

	server := extractExistingServerFromArg(p.db, event.Discord(), event.Fields()[2])
	if server == nil {
		event.Respond("serverlist.hide.no-server") // nolint: errcheck
		return
	}

	if !stringSliceContains(event.UserID, server.EditorUserIDs) {
		event.Respond("serverlist.hide.no-editor") // nolint: errcheck
		return
	}

	if server.State != StatePublic && server.State != StateHidden {
		event.Respond("serverlist.hide.wrong-status") // nolint: errcheck
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
