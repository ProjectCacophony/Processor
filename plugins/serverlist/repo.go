package serverlist

import (
	"github.com/jinzhu/gorm"
)

func categoryCreate(
	db *gorm.DB,
	keywords []string,
	botID string,
	guildID string,
	channelID string,
	addedBy string,
	sortBy []SortBy,
	groupBy GroupBy,
) error {
	sortByValue := make([]string, len(sortBy))
	for i := range sortBy {
		sortByValue[i] = string(sortBy[i])
	}

	return db.Create(&Category{
		Keywords:  keywords,
		BotID:     botID,
		GuildID:   guildID,
		ChannelID: channelID,
		AddedBy:   addedBy,
		SortBy:    sortByValue,
		GroupBy:   groupBy,
	}).Error
}

func categoriesFindMany(db *gorm.DB, where ...interface{}) ([]Category, error) {
	var entries []Category

	err := db.Find(&entries, where...).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return entries, err
}
