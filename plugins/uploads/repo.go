package uploads

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/events"
)

func addUpload(
	db *gorm.DB,
	fileInfo *events.FileInfo,
	userID string,
) error {
	if fileInfo == nil {
		return errors.New("invalid file info passed")
	}

	return db.Create(&Upload{
		FileInfoID: fileInfo.ID,
		UserID:     userID,
	}).Error
}
