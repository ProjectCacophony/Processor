package lastfm

import (
	"github.com/jinzhu/gorm"
)

func getLastFmUsername(db *gorm.DB, userID string) string {
	var user User
	db.Where(&User{UserID: userID}).First(&user)

	return user.LastFMUsername
}

func setLastFmUsername(db *gorm.DB, userID, lastFmUsername string) error {
	var user User

	err := db.Where(&User{
		UserID: userID,
	}).Assign(&User{
		LastFMUsername: lastFmUsername,
	}).FirstOrCreate(&user).Error
	if err != nil {
		return err
	}

	return nil
}
