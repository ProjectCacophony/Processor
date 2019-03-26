package serverlist

import (
	"errors"
	"strings"
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

func (g *GroupBy) ChannelName(serverName string) string {
	if *g == GroupByAlphabet {
		serverName = strings.ToLower(serverName)

		if len(serverName) > 0 && serverNameInitialRegexp.MatchString(string(serverName[0])) {
			return string(serverName[0]) + "—"
		}
	}

	return "etc"
}

// nolint: gochecknoglobals
var allGroupBys = []GroupBy{
	GroupByAlphabet,
}

type State string

const (
	StateQueued   State = "queued"
	StatePublic   State = "public"
	StateExpired  State = "expired" // TODO
	StateRemoved  State = "removed" // TODO
	StateRejected State = "rejected"
	StateHidden   State = "hidden" // TODO
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
}

func (*Server) TableName() string {
	return "serverlist_servers"
}

func (s *Server) QueueApprove(p *Plugin, guildID string) error {
	if s == nil {
		return errors.New("server is nil")
	}

	err := serverSetState(p.db, s.ID, StatePublic)
	if err != nil {
		return err
	}

	p.refreshQueue(guildID)
	return nil
}

func (s *Server) QueueReject(p *Plugin, guildID string) error {
	if s == nil {
		return errors.New("server is nil")
	}

	err := serverSetState(p.db, s.ID, StateRejected)
	if err != nil {
		return err
	}

	p.refreshQueue(guildID)
	return nil
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
