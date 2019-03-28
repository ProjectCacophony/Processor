package serverlist

import (
	"errors"
	"time"

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

func serverAdd(
	db *gorm.DB,
	names []string,
	description string,
	inviteCode string,
	guildID string,
	editorUserIDs []string,
	categoryIDs []uint,
	totalMembers int,
	botID string,
) error {
	server := &Server{
		Names:         names,
		Description:   description,
		InviteCode:    inviteCode,
		GuildID:       guildID,
		EditorUserIDs: editorUserIDs,
		TotalMembers:  totalMembers,
		State:         StateQueued,
		LastChecked:   time.Now(),
		BotID:         botID,
	}

	err := db.Create(server).Error
	if err != nil {
		return err
	}

	for _, categoryID := range categoryIDs {
		err := db.Create(&ServerCategory{
			ServerID:   server.ID,
			CategoryID: categoryID,
		}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func serverFind(db *gorm.DB, query string, where ...interface{}) (*Server, error) {
	var entry Server

	err := db.
		Preload("Categories.Category").
		Where(query, where...).
		First(&entry).
		Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func serversFindMany(db *gorm.DB, where ...interface{}) ([]*Server, error) {
	var entries []*Server

	err := db.
		Preload("Categories.Category").
		Find(&entries, where...).
		Order("created_at DESC").
		Error
	if err != nil {
		return nil, err
	}
	return entries, err
}

func serverRemove(db *gorm.DB, id uint) error {
	if id == 0 {
		return errors.New("please specify which server to delete")
	}

	return db.Delete(Server{}, "id = ?", id).Error
}

func serverSetState(db *gorm.DB, id uint, state State) error {
	return serverSetStateWithReason(db, id, state, "")
}

func serverSetStateWithReason(db *gorm.DB, id uint, state State, reason string) error {
	return db.Model(Server{}).Where("id = ?", id).Update(Server{
		State:  state,
		Reason: reason,
	}).Error
}
