package customcommands

import (
	"errors"

	"github.com/jinzhu/gorm"
)

const (
	noContent string = "No content or attachement."
)

func createCustomCommand(db *gorm.DB, command CustomCommand) error {
	if command.Content == "" && command.File == nil {
		return errors.New(noContent)
	}

	err := db.Create(&command).Error
	if err != nil {
		return err
	}

	if command.File != nil && command.Model.ID != 0 {

		command.File.CustomCommandID = command.Model.ID
		err := db.Save(&command.File).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func upsertCustomCommand(db *gorm.DB, command *CustomCommand) error {
	if command == nil {
		return errors.New("command cannot be nil")
	}

	if command.File != nil && command.Model.ID != 0 {

		command.File.CustomCommandID = command.Model.ID
		err := db.Save(&command.File).Error
		if err != nil {
			return err
		}
	}

	return db.Save(command).Error
}

func removeCustomCommand(db *gorm.DB, command *CustomCommand) error {
	if command.Model.ID == 0 {
		return errors.New("empty ID passed")
	}

	return db.Delete(command).Error
}
