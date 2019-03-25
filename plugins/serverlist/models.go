package serverlist

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type SortBy string

const (
	SortByAlphabetical SortBy = "alphabetical"
	SortByAddingDate   SortBy = "adding_date"
	SortByMemberCount  SortBy = "member_count"
)

// nolint: gochecknoglobals
var allSortBys = []SortBy{
	SortByAlphabetical,
	SortByAddingDate,
	SortByMemberCount,
}

type GroupBy string

const (
	GroupByAlphabet GroupBy = "alphabetical"
)

// nolint: gochecknoglobals
var allGroupBys = []GroupBy{
	GroupByAlphabet,
}

type Category struct {
	gorm.Model
	Keywords  pq.StringArray `gorm:"type:varchar[]"`
	BotID     string
	GuildID   string
	ChannelID string
	AddedBy   string
	SortBy    pq.StringArray `gorm:"type:varchar[]"`
	GroupBy   GroupBy
}

func (*Category) TableName() string {
	return "serverlist_categories"
}
