package serverlist

import (
	"errors"
	"strings"
	"time"

	"gitlab.com/Cacophony/go-kit/discord"

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
	StateRejected State = "rejected"
	StateHidden   State = "hidden"   // TODO
	StateCensored State = "censored" // TODO
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
	Reason        string
	LastChecked   time.Time
	BotID         string
}

func (*Server) TableName() string {
	return "serverlist_servers"
}

func (s *Server) QueueApprove(p *Plugin) error {
	if s == nil {
		return errors.New("server is nil")
	}

	if s.State != StateQueued {
		return errors.New("can only approve servers that are queued")
	}

	err := serverSetState(p.db, s.ID, StatePublic)
	if err != nil {
		return err
	}

	for _, category := range s.Categories {
		p.refreshQueue(category.Category.GuildID)
	}

	session, err := discord.NewSession(p.tokens, s.BotID)
	if err != nil {
		return err
	}

	for _, editorUserID := range s.EditorUserIDs {
		discord.SendComplexWithVars( // nolint: errcheck
			p.redis,
			session,
			p.Localisations(),
			editorUserID,
			&discordgo.MessageSend{
				Content: "serverlist.dm.server-approved",
			},
			true,
			"server",
			s,
		)
	}

	return nil
}

func (s *Server) QueueReject(p *Plugin, reason string) error {
	if s == nil {
		return errors.New("server is nil")
	}

	if s.State != StateQueued {
		return errors.New("can only approve servers that are queued")
	}

	err := serverSetStateWithReason(p.db, s.ID, StateRejected, reason)
	if err != nil {
		return err
	}

	for _, category := range s.Categories {
		p.refreshQueue(category.Category.GuildID)
	}

	session, err := discord.NewSession(p.tokens, s.BotID)
	if err != nil {
		return err
	}

	for _, editorUserID := range s.EditorUserIDs {
		discord.SendComplexWithVars( // nolint: errcheck
			p.redis,
			session,
			p.Localisations(),
			editorUserID,
			&discordgo.MessageSend{
				Content: "serverlist.dm.server-rejected",
			},
			true,
			"server",
			s,
			"reason",
			reason,
		)
	}

	return nil
}

func (s *Server) Remove(p *Plugin) error {
	if s == nil {
		return errors.New("server is nil")
	}

	err := serverRemove(p.db, s.ID)
	if err != nil {
		return err
	}

	for _, category := range s.Categories {
		p.refreshQueue(category.Category.GuildID)
	}

	session, err := discord.NewSession(p.tokens, s.BotID)
	if err != nil {
		return err
	}

	for _, editorUserID := range s.EditorUserIDs {
		discord.SendComplexWithVars( // nolint: errcheck
			p.redis,
			session,
			p.Localisations(),
			editorUserID,
			&discordgo.MessageSend{
				Content: "serverlist.dm.server-removed",
			},
			true,
			"server",
			s,
		)
	}

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
