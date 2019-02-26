package whitelist

import (
	"github.com/jinzhu/gorm"
)

func whitelistAddServer(db *gorm.DB, userID, guildID string) error {
	return db.Create(&Entry{
		WhitelistedByUserID: userID,
		GuildID:             guildID,
	}).Error
}

func whitelistRemoveServer(db *gorm.DB, guildID string) error {
	return db.Delete(Entry{}, "guild_id = ?", guildID).Error
}

func whitelistGetAllServers(db *gorm.DB) ([]Entry, error) {
	var entries []Entry

	err := db.Find(&entries).Order("created_at DESC").Error
	return entries, err
}

func blacklistAddServer(db *gorm.DB, userID, guildID string) error {
	return db.Create(&BlacklistEntry{
		BlacklistedByUserID: userID,
		GuildID:             guildID,
	}).Error
}

func blacklistFindServer(db *gorm.DB, guildID string) (*BlacklistEntry, error) {
	var entry BlacklistEntry

	err := db.First(&entry, "guild_id = ?", guildID).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func blacklistGetAllServers(db *gorm.DB) ([]BlacklistEntry, error) {
	var entries []BlacklistEntry

	err := db.Find(&entries).Order("created_at DESC").Error
	return entries, err
}
