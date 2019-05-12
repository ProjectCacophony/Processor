package rss

import (
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"gitlab.com/Cacophony/go-kit/permissions"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/state"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localization"
	"go.uber.org/zap"
)

type Plugin struct {
	logger     *zap.Logger
	state      *state.State
	db         *gorm.DB
	parser     *gofeed.Parser
	httpClient *http.Client
}

func (p *Plugin) Name() string {
	return "rss"
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.state = params.State
	p.logger = params.Logger
	p.db = params.DB
	p.parser = gofeed.NewParser()
	p.httpClient = &http.Client{
		Timeout: time.Second * 30,
	}

	return params.DB.AutoMigrate(Entry{}).Error
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

func (p *Plugin) Localizations() []interfaces.Localization {
	local, err := localization.NewFileSource("assets/translations/rss.en.toml", "en")
	if err != nil {
		p.logger.Error("failed to load localization", zap.Error(err))
	}

	return []interfaces.Localization{local}
}

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name:        p.Name(),
		Description: "help.rss.description",
		ParamSets: []common.ParamSet{{
			Description:               "Add RSS feed to this channel or another channel if specified.",
			DiscordPermissionRequired: permissions.DiscordManageChannels,
			Params: []common.PluginParams{
				{Name: "add", Type: common.Text, NotVariable: true},
				{Name: "channel", Type: common.Channel, Optional: true},
				{Name: "link", Type: common.Text},
			},
		}, {
			Description:               "Remove RSS feed from this channel or another channel if specified.",
			DiscordPermissionRequired: permissions.DiscordManageChannels,
			Params: []common.PluginParams{
				{Name: "remove", Type: common.Text, NotVariable: true},
				{Name: "channel", Type: common.Channel, Optional: true},
				{Name: "link", Type: common.Text},
			},
		}},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "rss" {
		return false
	}

	if len(event.Fields()) > 1 {
		switch strings.ToLower(event.Fields()[1]) {
		case "add":
			event.RequireOr(func() {
				p.add(event)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true

		case "remove":
			event.RequireOr(func() {
				p.remove(event)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true
		}
	}

	p.status(event)

	return true
}
