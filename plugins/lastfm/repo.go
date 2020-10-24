package lastfm

import (
	"errors"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/config"
)

const configKey = "cacophony:lastfm:username"

func getLastFmUsername(db *gorm.DB, userID string) string {
	lastFMUsername, _ := config.UserGetString(db, userID, configKey)
	if lastFMUsername != "" {
		return lastFMUsername
	}

	// fallback to legacy storage
	var user User
	db.Where(&User{UserID: userID}).First(&user)

	return user.LastFMUsername
}

func setLastFmUsername(db *gorm.DB, userID, lastFmUsername string) error {
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}

	return config.UserSetString(db, userID, configKey, lastFmUsername)
}
