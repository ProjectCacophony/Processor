package whitelist

import (
	"github.com/jinzhu/gorm"
)

func whitelistAdd(db *gorm.DB, userID, guildID string) error {
	return db.Create(&Entry{
		WhitelistedByUserID: userID,
		GuildID:             guildID,
	}).Error
}

func whitelistRemove(db *gorm.DB, guildID string) error {
	return db.Delete(Entry{}, "guild_id = ?", guildID).Error
}

func whitelistFind(db *gorm.DB, where ...interface{}) (*Entry, error) {
	var entry Entry

	err := db.First(&entry, where...).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func whitelistFindMany(db *gorm.DB, where ...interface{}) ([]Entry, error) {
	var entries []Entry

	err := db.First(&entries, where...).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func whitelistAll(db *gorm.DB) ([]Entry, error) {
	var entries []Entry

	err := db.Find(&entries).Order("created_at DESC").Error
	return entries, err
}

func blacklistAdd(db *gorm.DB, userID, guildID string) error {
	return db.Create(&BlacklistEntry{
		BlacklistedByUserID: userID,
		GuildID:             guildID,
	}).Error
}

// func blacklistRemove(db *gorm.DB, guildID string) error {
// 	return db.Delete(BlacklistEntry{}, "guild_id = ?", guildID).Error
// }

func blacklistFind(db *gorm.DB, where ...interface{}) (*BlacklistEntry, error) {
	var entry BlacklistEntry

	err := db.First(&entry, where...).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func blacklistAll(db *gorm.DB) ([]BlacklistEntry, error) {
	var entries []BlacklistEntry

	err := db.Find(&entries).Order("created_at DESC").Error
	return entries, err
}
