package customcommands

import (
	"time"

	"github.com/jinzhu/gorm"
)

func entryAdd(
	db *gorm.DB,
	name string,
	userID string,
	guildID string,
	content string,
	isUserCommand bool,
) error {
	return db.Create(&Entry{
		Name:          name,
		UserID:        userID,
		GuildID:       guildID,
		Content:       content,
		Date:          time.Now(),
		IsUserCommand: isUserCommand,
	}).Error
}

func entryFind(db *gorm.DB, where ...interface{}) (*Entry, error) {
	var entry Entry

	err := db.First(&entry, where...).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return &entry, err
}

func entryFindMany(db *gorm.DB, where ...interface{}) ([]Entry, error) {
	var entries []Entry

	err := db.Find(&entries, where...).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return entries, err
}
