package uploads

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleUpload(event *events.Event) {
	// TODO: allow multiple uploads at once

	var uploads []Upload

	for _, attachment := range event.MessageCreate.Attachments {
		file, err := event.AddAttachement(attachment)
		if err != nil {
			event.Except(err)
			return
		}

		err = addUpload(p.db, file, event.UserID)
		if err != nil {
			event.Except(err)
			return
		}

		uploads = append(uploads, Upload{
			FileInfo: *file,
			UserID:   event.UserID,
		})
	}

	if len(uploads) <= 0 {
		event.Respond("uploads.upload.too-few")
		return
	}

	_, err := event.Respond("uploads.upload.content", "uploads", uploads)
	event.Except(err)
}
