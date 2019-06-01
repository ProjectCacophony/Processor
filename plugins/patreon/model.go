package patreon

import (
	"github.com/jinzhu/gorm"
)

// model maintained by Worker
type Patron struct {
	gorm.Model
	PatreonUserID string

	FirstName     string
	VanityName    string
	FullName      string
	PatronStatus  string
	DiscordUserID string
}

func (*Patron) TableName() string {
	return "patrons"
}
