package customcommands

import (
	"errors"
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

func entryUpsert(db *gorm.DB, entry *Entry) error {
	if entry == nil {
		return errors.New("entry cannot be nil")
	}

	err := db.
		Where("id = ?", entry.Model.ID).
		Assign(entry).
		FirstOrCreate(&Entry{}).
		Error
	if err != nil {
		return err
	}

	return nil
}

func entryRemove(db *gorm.DB, id uint) error {
	if id == 0 {
		return errors.New("empty ID passed")
	}

	return db.Delete(Entry{}, "id = ?", id).Error
}

func entryUpdateTriggered(db *gorm.DB, entry *Entry) error {
	if entry == nil {
		return errors.New("entry cannot be nil")
	}
	return db.Model(entry).Update("triggered", gorm.Expr("triggered + 1")).Error
}
