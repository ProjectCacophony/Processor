package eventlog

import (
	"context"
	"errors"

	"github.com/getsentry/raven-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func CreateItem(db *gorm.DB, publisher *events.Publisher, item *Item) error {
	if item == nil {
		return errors.New("item cannot be empty")
	}
	if item.GuildID == "" {
		return errors.New("GuildID cannot be empty")
	}

	// decorate item
	item.UUID = uuid.New()

	// prepare event
	event, err := events.New(events.CacophonyEventlogUpdate)
	if err != nil {
		return err
	}
	event.EventlogUpdate = &events.EventlogUpdate{
		GuildID: item.GuildID,
	}
	event.GuildID = item.GuildID

	err = db.Create(&item).Error
	if err != nil {
		return err
	}

	event.EventlogUpdate.ItemID = item.ID

	err, recoverable := publisher.Publish(context.Background(), event)
	if err != nil && !recoverable {
		raven.CaptureError(err, nil)
		zap.L().Fatal(
			"received unrecoverable error while publishing \"sudo as\" message",
			zap.Error(err),
		)
	}
	return err
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

func saveItemMessage(db *gorm.DB, id uint, messageID, channelID string) error {
	return db.Model(&Item{}).Where("id = ?", id).Updates(Item{
		LogMessage: ItemLogMessage{
			MessageID: messageID,
			ChannelID: channelID,
		},
	}).Error
}
