package rss

import (
	"errors"

	"github.com/jinzhu/gorm"
)

func entryAdd(
	db *gorm.DB,
	userID string,
	channelID string,
	guildID string,
	name string,
	url string,
	feedURL string,
	botID string,
	dm bool,
) error {
	return db.Create(&Entry{
		Model:     gorm.Model{},
		GuildID:   guildID,
		ChannelID: channelID,
		AddedBy:   userID,
		BotID:     botID,
		DM:        dm,
		Name:      name,
		URL:       url,
		FeedURL:   feedURL,
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

func entryRemove(db *gorm.DB, id uint) error {
	if id == 0 {
		return errors.New("empty ID passed")
	}

	return db.Delete(Entry{}, "id = ?", id).Error
}

func countPosts(db *gorm.DB, id uint) (int, error) {
	var amount int
	err := db.Model(&Post{}).Where("entry_id = ?", id).Count(&amount).Error
	return amount, err
}
