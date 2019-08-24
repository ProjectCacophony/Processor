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

func getUploads(
	db *gorm.DB,
	userID string,
) ([]Upload, error) {
	var uploads []Upload

	err := db.
		Preload("FileInfo").
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&uploads).Error

	return uploads, err
}

func deleteUpload(db *gorm.DB, id uint) error {
	if id == 0 {
		return errors.New("empty ID passed")
	}

	return db.Delete(Upload{}, "id = ?", id).Error
}
