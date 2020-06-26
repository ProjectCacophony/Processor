package roles

import "gitlab.com/Cacophony/go-kit/events"

func (p *Plugin) createAutoRole(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	// TODO: more than 5 params has duration
	serverRoleID := event.Fields()[3]
	if serverRoleID == "" {
		event.Respond("roles.role.role-not-found-on-server")
		return
	}

	event.Respond("Not implemented cause i'm lazy")
}
