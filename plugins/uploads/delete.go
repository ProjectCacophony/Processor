package uploads

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleDelete(event *events.Event) {
	uploads, err := getUploads(event.DB(), event.UserID)
	if err != nil {
		event.Except(err)
		return
	}

	for _, upload := range uploads {
		for _, field := range event.Fields() {
			if upload.FileInfo.GetLink() != field {
				continue
			}

			err = event.DeleteFile(&upload.FileInfo)
			if err != nil {
				event.Except(err)
				return
			}

			err = deleteUpload(p.db, upload.ID)
			if err != nil {
				event.Except(err)
				return
			}

			_, err = event.Respond("uploads.delete.success", "upload", upload)
			event.Except(err)
			return
		}
	}

	event.Respond("uploads.delete.not-found")
	return
}
