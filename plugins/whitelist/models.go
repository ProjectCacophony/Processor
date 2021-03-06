package whitelist

import (
	"github.com/jinzhu/gorm"
)

type Entry struct {
	gorm.Model
	WhitelistedByUserID string
	GuildID             string
}

func (*Entry) TableName() string {
	return "whitelist_entries"
}

type BlacklistEntry struct {
	gorm.Model
	BlacklistedByUserID string
	GuildID             string
}

func (*BlacklistEntry) TableName() string {
	return "whitelist_blacklist_entries"
}
