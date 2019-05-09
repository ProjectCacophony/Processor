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

var allGroupBys = []GroupBy{
	GroupByAlphabet,
}

type State string

const (
	StateQueued   State = "queued"
	StatePublic   State = "public"
	StateExpired  State = "expired"
	StateHidden   State = "hidden"
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
	Change        ServerChange `gorm:"embedded;embedded_prefix:change_"`
}

func (*Server) TableName() string {
	return "serverlist_servers"
}

type ServerChange struct {
	Names       pq.StringArray `gorm:"type:varchar[]"`
	Description string
	InviteCode  string
	Categories  pq.Int64Array `gorm:"type:integer[]"` // TODO: uint array
	State       State
}

func (s *Server) refresh(p *Plugin) error {
	server, err := serverFind(p.db, "id = ?", s.ID)
	if err != nil {
		return err
	}

	*s = *server
	return nil
}

func (s *Server) QueueApprove(p *Plugin) error {
	if s == nil {
		return errors.New("server is nil")
	}

	if s.State != StateQueued {
		return errors.New("can only approve servers that are queued")
	}

	newState := StatePublic

	if s.Change.State != "" {
		err := s.ApplyChange(p, s.Change)
		if err != nil {
			return err
		}

		if s.Change.State != StateQueued {
			newState = s.Change.State
		}
	}

	err := serverSetState(p.db, s.ID, newState)
	if err != nil {
		return err
	}

	err = s.refresh(p)
	if err != nil {
		return err
	}

	s.refreshQueues(p)

	p.refreshList(s.BotID)

	session, err := discord.NewSession(p.tokens, s.BotID)
	if err != nil {
		return err
	}

	for _, editorUserID := range s.EditorUserIDs {
		discord.SendComplexWithVars(
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

	return p.sendLogMessageForServer(
		session,
		s,
		&discordgo.MessageSend{
			Embed: p.getLogApprovedEmbed(s),
		},
	)
}

func (s *Server) QueueReject(p *Plugin, reason string) error {
	if s == nil {
		return errors.New("server is nil")
	}

	if s.State != StateQueued {
		return errors.New("can only reject servers that are queued")
	}

	var err error

	if s.Change.State != "" && s.Change.State != StateQueued {
		// rejecting change, reverting server back to previous state

		err = s.resetChange(p)
		if err != nil {
			return err
		}

		err = serverSetStateWithReason(p.db, s.ID, s.Change.State, reason)
		if err != nil {
			return err
		}

		s.State = s.Change.State

	} else {
		// rejecting server, removing server
		err = s.Remove(p, false)
		if err != nil {
			return err
		}
	}

	s.refreshQueues(p)

	if s.State == StatePublic {
		p.refreshList(s.BotID)
	}

	session, err := discord.NewSession(p.tokens, s.BotID)
	if err != nil {
		return err
	}

	for _, editorUserID := range s.EditorUserIDs {
		discord.SendComplexWithVars(
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

func (s *Server) Remove(p *Plugin, editor bool) error {
	if s == nil {
		return errors.New("server is nil")
	}

	if s.State == StateCensored {
		return errors.New("can not remove servers that are censored")
	}

	err := serverRemove(p.db, s.ID)
	if err != nil {
		return err
	}

	s.refreshQueues(p)

	p.refreshList(s.BotID)

	session, err := discord.NewSession(p.tokens, s.BotID)
	if err != nil {
		return err
	}

	for _, editorUserID := range s.EditorUserIDs {
		discord.SendComplexWithVars(
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
			"editor",
			editor,
		)
	}

	return nil
}

func (s *Server) Hide(p *Plugin) error {
	if s == nil {
		return errors.New("server is nil")
	}

	if s.State != StatePublic {
		return errors.New("can only hide servers that are public")
	}

	err := serverSetState(p.db, s.ID, StateHidden)
	if err != nil {
		return err
	}

	s.refreshQueues(p)

	p.refreshList(s.BotID)

	session, err := discord.NewSession(p.tokens, s.BotID)
	if err != nil {
		return err
	}

	for _, editorUserID := range s.EditorUserIDs {
		discord.SendComplexWithVars(
			p.redis,
			session,
			p.Localisations(),
			editorUserID,
			&discordgo.MessageSend{
				Content: "serverlist.dm.server-hidden",
			},
			true,
			"server",
			s,
		)
	}

	return nil
}

func (s *Server) Unhide(p *Plugin) error {
	if s == nil {
		return errors.New("server is nil")
	}

	if s.State != StateHidden {
		return errors.New("can only unhide servers that are hidden")
	}

	err := serverSetState(p.db, s.ID, StatePublic)
	if err != nil {
		return err
	}

	s.refreshQueues(p)

	p.refreshList(s.BotID)

	session, err := discord.NewSession(p.tokens, s.BotID)
	if err != nil {
		return err
	}

	for _, editorUserID := range s.EditorUserIDs {
		discord.SendComplexWithVars(
			p.redis,
			session,
			p.Localisations(),
			editorUserID,
			&discordgo.MessageSend{
				Content: "serverlist.dm.server-unhidden",
			},
			true,
			"server",
			s,
		)
	}

	return nil
}

func (s *Server) Edit(p *Plugin, changes ServerChange) error {
	if s == nil {
		return errors.New("server is nil")
	}

	if s.State == StateCensored {
		return errors.New("can not remove servers that are censored")
	}

	var change Server

	if s.Change.State == "" {
		change.Change.State = s.State
	}

	if len(changes.Names) > 0 {
		change.Change.Names = changes.Names
	}
	if len(changes.Description) > 0 {
		change.Change.Description = changes.Description
	}
	if len(changes.InviteCode) > 0 {
		change.Change.InviteCode = changes.InviteCode
	}
	if len(changes.Categories) > 0 {
		change.Change.Categories = changes.Categories
	}

	if s.State == StateQueued {
		err := s.ApplyChange(p, changes)
		if err != nil {
			return err
		}
	} else {
		err := serverUpdate(p.db, s.ID, change)
		if err != nil {
			return err
		}
	}

	if s.State != StateQueued {
		err := serverSetState(p.db, s.ID, StateQueued)
		if err != nil {
			return err
		}
	}

	s.refreshQueues(p)

	p.refreshList(s.BotID)

	return nil
}

func (s *Server) ApplyChange(p *Plugin, change ServerChange) error {
	serverCategories := make([]ServerCategory, len(change.Categories))

	for i, categoryID := range change.Categories {
		serverCategories[i] = ServerCategory{
			ServerID:   s.ID,
			CategoryID: uint(categoryID),
		}
	}

	update := Server{
		Names:       change.Names,
		Description: change.Description,
		InviteCode:  change.InviteCode,
		Categories:  serverCategories,
	}

	err := serverUpdate(p.db, s.ID, update)
	if err != nil {
		return err
	}

	return s.resetChange(p)
}

func (s *Server) Update(p *Plugin, update Server) error {
	err := serverUpdate(p.db, s.ID, update)
	if err != nil {
		return err
	}

	s.refreshQueues(p)

	p.refreshList(s.BotID)

	return nil
}

func (s *Server) resetChange(p *Plugin) error {
	return serverResetChange(p.db, s.ID)
}

func (s *Server) refreshQueues(p *Plugin) {
	triedGuildIDs := make(map[string]interface{})

	for _, category := range s.Categories {
		if _, ok := triedGuildIDs[category.Category.GuildID]; ok {
			continue
		}

		triedGuildIDs[category.Category.GuildID] = nil

		p.refreshQueue(category.Category.GuildID)
	}
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
