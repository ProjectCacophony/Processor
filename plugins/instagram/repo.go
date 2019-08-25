package instagram

import (
	"errors"

	"github.com/jinzhu/gorm"
)

func entryAdd(
	db *gorm.DB,
	userID string,
	channelOrUserID string,
	guildID string,
	instagramUsername string,
	instagramAccountID string,
	botID string,
	dm bool,
) error {
	return db.Create(&Entry{
		Model:              gorm.Model{},
		GuildID:            guildID,
		ChannelOrUserID:    channelOrUserID,
		DM:                 dm,
		AddedBy:            userID,
		BotID:              botID,
		InstagramUsername:  instagramUsername,
		InstagramAccountID: instagramAccountID,
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

type modifyType int

const (
	modifyPosts modifyType = iota
	modifyStory
)

func entryModify(db *gorm.DB, id uint, modification modifyType, value bool) error {
	updates := &Entry{}
	switch modification {
	case modifyPosts:
		updates.DisablePostFeed = value
	case modifyStory:
		updates.DisableStoryFeed = value
	}

	return db.Model(&Entry{}).Where("id = ?", id).Updates(updates).Error
}
