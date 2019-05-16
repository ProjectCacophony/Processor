package serverlist

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin struct {
	logger        *zap.Logger
	db            *gorm.DB
	state         *state.State
	redis         *redis.Client
	tokens        map[string]string
	staffRoles    interfaces.Permission
	localizations []interfaces.Localization
}

func (p *Plugin) Name() string {
	return "serverlist"
}

func (p *Plugin) Start(params common.StartParameters) error {
	var err error

	p.logger = params.Logger
	p.db = params.DB
	p.state = params.State
	p.redis = params.Redis
	p.tokens = params.Tokens
	p.localizations = params.Localizations

	p.staffRoles = permissions.Or(
		permissions.BotAdmin,
		// Test / Staff
		permissions.NewDiscordRole(p.state, "561619599129444390", "561619665197989893"),
	)

	err = p.db.AutoMigrate(
		Category{},
		Server{},
		ServerCategory{},
	).Error
	return err
}

func (p *Plugin) Stop(params common.StopParameters) error {
	return nil
}

func (p *Plugin) Priority() int {
	return 0
}

func (p *Plugin) Passthrough() bool {
	return false
}

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name:        p.Name(),
		Description: "serverlist.help.description",
		Commands: []common.Command{{
			Name: "List all Categories",
			Params: []common.CommandParam{
				{Name: "category", Type: common.Hardcoded},
			},
		}, {
			Name:                "Create a Category",
			Description:         "Sort By Options: `alphabetical`, `adding_date`, `member_count`\n\t\tGroup By Options: `alphabetical`",
			PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
			Params: []common.CommandParam{
				{Name: "category", Type: common.Hardcoded},
				{Name: "create", Type: common.Hardcoded},
				{Name: "channel", Type: common.Channel},
				{Name: "keywords, seporate by ;", Type: common.QuotedText},
				{Name: "sory by, seporate by ;", Type: common.QuotedText},
				{Name: "group by, for categories", Type: common.QuotedText, Optional: true},
			},
		}, {
			Name:                "Set the Serverlist Log Channel",
			PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
			Params: []common.CommandParam{
				{Name: "log", Type: common.Hardcoded},
			},
		}, {
			Name:                "Set the Serverlist Queue Channel",
			PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
			Params: []common.CommandParam{
				{Name: "queue", Type: common.Hardcoded},
			},
		}, {
			Name:                "Reject the current Server",
			PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
			Params: []common.CommandParam{
				{Name: "reject", Type: common.Hardcoded},
				{Name: "reason", Type: common.QuotedText},
			},
		}, {
			Name:                "Censor a Server",
			Description:         "Run the command again to uncensor a Server.",
			PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
			Params: []common.CommandParam{
				{Name: "censor", Type: common.Hardcoded},
				{Name: "Server Invite", Type: common.Text},
				{Name: "reason", Type: common.QuotedText},
			},
		}, {
			Name:                "Manually refresh the Queue Channel",
			PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
			Params: []common.CommandParam{
				{Name: "queue", Type: common.Hardcoded},
				{Name: "refresh", Type: common.Hardcoded},
			},
		}, {
			Name:                "Manually refresh the List Channels",
			PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
			Params: []common.CommandParam{
				{Name: "list", Type: common.Hardcoded},
				{Name: "refresh", Type: common.Hardcoded},
			},
		}, {
			Name:                "Clear the List Channel Cache",
			PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
			Params: []common.CommandParam{
				{Name: "list", Type: common.Hardcoded},
				{Name: "clear-cache", Type: common.Hardcoded},
			},
		}, {
			Name: "Submit a Server",
			Params: []common.CommandParam{
				{Name: "add", Type: common.Hardcoded},
				{Name: "Server Invite", Type: common.Text},
				{Name: "Description", Type: common.QuotedText},
				{Name: "Category, if multiple, separated by ;", Type: common.QuotedText},
			},
		}, {
			Name: "List your Servers",
			Params: []common.CommandParam{
				{Name: "list", Type: common.Hardcoded},
			},
		}, {
			Name:        "Hide a Servers",
			Description: "Run the command again to unhide the Server.",
			Params: []common.CommandParam{
				{Name: "hide", Type: common.Hardcoded},
				{Name: "Server Invite", Type: common.Text},
			},
		}, {
			Name:        "Edit a Server Name",
			Description: "When editing the Server Name, the Server goes back to the Queue as the change has to be approved by a Moderator. \n\t\tIn lists that are sorted by \"Server Added At\" the position will not be lost after approval, as the Adding Date will not be changed.",
			Params: []common.CommandParam{
				{Name: "edit", Type: common.Hardcoded},
				{Name: "Server Invite", Type: common.Text},
				{Name: "name", Type: common.Hardcoded},
				{Name: "New Server Name, if multiple, separated by ;", Type: common.QuotedText},
			},
		}, {
			Name:        "Edit a Server Description",
			Description: "When editing the Server Description, the Server goes back to the Queue as the change has to be approved by a Moderator. \n\t\tIn lists that are sorted by \"Server Added At\" the position will not be lost after approval, as the Adding Date will not be changed.",
			Params: []common.CommandParam{
				{Name: "edit", Type: common.Hardcoded},
				{Name: "Server Invite", Type: common.Text},
				{Name: "description", Type: common.Hardcoded},
				{Name: "New Server Description", Type: common.QuotedText},
			},
		}, {
			Name:        "Edit a Server Category",
			Description: "When editing the Server Category, the Server goes back to the Queue as the change has to be approved by a Moderator. \n\t\tIn lists that are sorted by \"Server Added At\" the position will not be lost after approval, as the Adding Date will not be changed.",
			Params: []common.CommandParam{
				{Name: "edit", Type: common.Hardcoded},
				{Name: "Server Invite", Type: common.Text},
				{Name: "category ", Type: common.Hardcoded},
				{Name: "New Category, if multiple, separated by ;", Type: common.QuotedText},
			},
		}, {
			Name:        "Edit a Server Invite",
			Description: "Adds as an Editor if the User is not an Editor so far.\n\t\tRemoves as an Editor if the User is an Editor.\n\t\tAll Editors have equal rights about the Listing, they can all edit, hide, and remove the Server.\n\t\tPlease make sure you only add Users you trust as Editors",
			Params: []common.CommandParam{
				{Name: "edit", Type: common.Hardcoded},
				{Name: "Server Invite", Type: common.Text},
				{Name: "invite", Type: common.Hardcoded},
				{Name: "New Server Invite", Type: common.QuotedText},
			},
		}, {
			Name: "Add/Remove Server Editors",
			Params: []common.CommandParam{
				{Name: "edit", Type: common.Hardcoded},
				{Name: "Server Invite", Type: common.Text},
				{Name: "editor", Type: common.Hardcoded},
				{Name: "user", Type: common.User},
			},
		}, {
			Name: "Remove a Server",
			Params: []common.CommandParam{
				{Name: "remove", Type: common.Hardcoded},
				{Name: "Server Invite", Type: common.Text},
			},
		}},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	switch event.Type {

	case events.MessageReactionAddType:
		return p.handleQueueReaction(event)

	case events.CacophonyServerlistServerExpire:
		p.handleExpired(event)
		return true

	}

	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "serverlist" {
		return false
	}

	if len(event.Fields()) >= 2 {
		switch event.Fields()[1] {
		case "category", "categories":

			if len(event.Fields()) >= 3 {
				switch event.Fields()[2] {
				case "create", "add":

					event.Require(func() {

						p.handleCategoryCreate(event)
					},
						permissions.Not(permissions.DiscordChannelDM),
						p.staffRoles,
					)
					return true
				}
			}

			p.handleCategoryStatus(event)
			return true

		case "add":

			p.handleAdd(event)
			return true

		case "queue":

			if len(event.Fields()) >= 3 {
				if event.Fields()[2] == "refresh" {

					event.Require(func() {

						p.handleQueueRefresh(event)
					},
						permissions.Not(permissions.DiscordChannelDM),
						p.staffRoles,
					)
					return true
				}
			}

			event.Require(func() {

				p.handleQueue(event)
			},
				permissions.Not(permissions.DiscordChannelDM),
				p.staffRoles,
			)
			return true

		case "list":

			if len(event.Fields()) >= 3 {

				switch event.Fields()[2] {

				case "refresh":

					event.Require(func() {

						p.handleListRefresh(event)
					},
						permissions.Not(permissions.DiscordChannelDM),
						p.staffRoles,
					)
					return true

				case "clear-cache":

					event.Require(func() {

						p.handleListClearCache(event)
					},
						permissions.Not(permissions.DiscordChannelDM),
						p.staffRoles,
					)
					return true

				}
			}

		case "reject":

			event.Require(func() {

				p.handleQueueReject(event)
			},
				permissions.Not(permissions.DiscordChannelDM),
				p.staffRoles,
			)
			return true

		case "remove":

			p.handleRemove(event)
			return true

		case "hide", "unhide":

			p.handleHide(event)
			return true

		case "log":

			event.Require(func() {

				p.handleLog(event)
			},
				permissions.Not(permissions.DiscordChannelDM),
				p.staffRoles,
			)
			return true

		case "edit":

			p.handleEdit(event)
			return true

		case "censor", "uncensor":

			event.Require(func() {

				p.handleCensor(event)
			},
				permissions.Not(permissions.DiscordChannelDM),
				p.staffRoles,
			)
			return true

		}
	}

	p.handleStatus(event)

	return true
}
