package serverlist

import (
	"time"

	"github.com/bwmarrin/discordgo"

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

type State string

const (
	StateQueued  State = "queued"
	StatePublic  State = "public"
	StateExpired State = "expired"
	StateRemoved State = "removed"
)

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

type Server struct {
	gorm.Model
	Names         pq.StringArray `gorm:"type:varchar[]"`
	Description   string
	InviteCode    string
	GuildID       string
	EditorUserIDs pq.StringArray `gorm:"type:varchar[]"`
	Categories    []ServerCategory
	TotalMembers  int
	State         State
	LastChecked   time.Time
	BotID         string
	// Messages      []Message
}

func (*Server) TableName() string {
	return "serverlist_servers"
}

type ServerCategory struct {
	gorm.Model
	ServerID   uint
	CategoryID uint
	Category   Category
}

func (*ServerCategory) TableName() string {
	return "serverlist_server_categories"
}

type QueueMessage struct {
	CurrentServerID uint
	MessageID       string
	Embed           *discordgo.MessageEmbed
}
