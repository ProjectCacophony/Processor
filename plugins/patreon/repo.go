package patreon

import (
	"github.com/jinzhu/gorm"
)

func getPatron(db *gorm.DB, userID string) (*Patron, error) {
	var patron Patron
	err := db.Where(&Patron{DiscordUserID: userID}).First(&patron).Error

	return &patron, err
}
