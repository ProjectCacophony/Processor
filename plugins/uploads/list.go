package uploads

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleList(event *events.Event) {
	uploads, err := getUploads(event.DB(), event.UserID)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("uploads.list.content", "uploads", uploads)
	event.Except(err)
}
