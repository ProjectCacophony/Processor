package greeter

import (
	"errors"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

func entryAdd(
	db *gorm.DB,
	guildID string,
	channelID string,
	greeterType greeterType,
	message string,
	rule *models.Rule,
) error {
	if rule == nil || rule.ID == 0 {
		return errors.New("invalid rule passed")
	}

	return db.Create(&Entry{
		GuildID:   guildID,
		ChannelID: channelID,
		Type:      greeterType,
		Message:   message,
		RuleID:    rule.ID,
	}).Error
}

func entriesFind(db *gorm.DB, guildID string) ([]Entry, error) {
	var entries []Entry

	err := db.
		Preload("Rule").
		Preload("Rule.Filters").
		Preload("Rule.Actions").
		Order("created_at DESC").
		Find(&entries, "guild_id = ?", guildID).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func entryFind(
	db *gorm.DB,
	guildID string,
	channelID string,
	greeterType greeterType,
) (*Entry, error) {
	var entry Entry

	err := db.
		Preload("Rule").
		Preload("Rule.Filters").
		Preload("Rule.Actions").
		Order("created_at DESC").
		First(
			&entry,
			"guild_id = ? AND channel_id = ? AND type = ?", guildID, channelID, greeterType,
		).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &entry, nil
}

func entryUpdate(
	db *gorm.DB,
	id uint,
	message string,
) error {
	return db.Model(Entry{}).Where("id = ?", id).Update("message", message).Error
}

func entryDelete(
	db *gorm.DB,
	id uint,
) error {
	if id == 0 {
		return errors.New("id cannot be empty")
	}

	return db.Where("id = ?", id).Delete(Entry{}).Error
}
