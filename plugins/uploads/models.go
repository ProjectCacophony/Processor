package uploads

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/events"
)

type Upload struct {
	gorm.Model
	FileInfoID uint
	FileInfo   events.FileInfo
	UserID     string
}

func (*Upload) TableName() string {
	return "uploads"
}
