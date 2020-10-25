package eventlog

import (
	"context"
	"errors"
	"strings"

	"github.com/getsentry/raven-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func isEnabled(event *events.Event) bool {
	enabled, err := config.GuildGetBool(event.DB(), event.GuildID, eventlogEnableKey)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		event.ExceptSilent(err)
		return false
	}

	return enabled
}

func CreateItem(db *gorm.DB, publisher *events.Publisher, item *Item) error {
	if item == nil {
		return errors.New("item cannot be empty")
	}
	if item.GuildID == "" {
		return errors.New("GuildID cannot be empty")
	}

	// decorate item
	item.UUID = uuid.New()

	err := db.Create(&item).Error
	if err != nil {
		return err
	}

	return publishUpdateEvent(publisher, item.GuildID, item.ID)
}

func CreateOptionForItem(db *gorm.DB, publisher *events.Publisher, id uint, guildID string, option *ItemOption) error {
	if id == 0 {
		return errors.New("id cannot be empty")
	}
	if option == nil {
		return errors.New("option cannot be empty")
	}
	option.ItemID = id

	err := db.
		Set("gorm:insert_option", "ON CONFLICT (\"item_id\", \"author_id\", \"key\") DO UPDATE SET \"updated_at\" = EXCLUDED.updated_at, \"previous_value\" = EXCLUDED.previous_value, \"new_value\" = EXCLUDED.new_value, \"type\" = EXCLUDED.type").
		Create(&option).Error
	if err != nil {
		return err
	}

	return publishUpdateEvent(publisher, guildID, id)
}

func GetItem(db *gorm.DB, id uint) (*Item, error) {
	if id == 0 {
		return nil, errors.New("id cannot be empty")
	}

	var item Item
	err := db.
		Preload("Options").
		Where("id = ?", id).
		First(&item).Error
	return &item, err
}

func FindItem(db *gorm.DB, where string, args ...interface{}) (*Item, error) {
	var item Item
	err := db.
		Preload("Options").
		Where(where, args...).
		First(&item).Error
	return &item, err
}

func FindManyItem(db *gorm.DB, limit int, where string, args ...interface{}) ([]Item, error) {
	var items []Item
	err := db.
		Preload("Options").
		Where(where, args...).
		Limit(limit).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

func saveItemMessage(db *gorm.DB, id uint, messageID, channelID string) error {
	if id == 0 {
		return errors.New("id cannot be empty")
	}

	return db.Model(&Item{}).Where("id = ?", id).Updates(Item{
		LogMessage: ItemLogMessage{
			MessageID: messageID,
			ChannelID: channelID,
		},
	}).Error
}

func markItemAsReverted(db *gorm.DB, publisher *events.Publisher, guildID string, id uint) error {
	if id == 0 {
		return errors.New("id cannot be empty")
	}

	err := db.Model(&Item{}).Where("id = ?", id).Updates(Item{
		Reverted: true,
	}).Error
	if err != nil {
		return err
	}

	return publishUpdateEvent(publisher, guildID, id)
}

func publishUpdateEvent(publisher *events.Publisher, guildID string, itemID uint) error {
	if publisher == nil {
		return nil
	}

	// prepare event
	event, err := events.New(events.CacophonyEventlogUpdate)
	if err != nil {
		return err
	}
	event.EventlogUpdate = &events.EventlogUpdate{
		GuildID: guildID,
		ItemID:  itemID,
	}
	event.GuildID = guildID

	err, recoverable := publisher.Publish(context.Background(), event)
	if err != nil && !recoverable {
		raven.CaptureError(err, nil)
		zap.L().Fatal(
			"received unrecoverable error while publishing \"create item\" message",
			zap.Error(err),
		)
	}
	return err
}
