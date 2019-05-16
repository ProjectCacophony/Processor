package instagram

import (
	"net/http"
	"strings"
	"time"

	"github.com/Seklfreak/ginsta"

	"gitlab.com/Cacophony/go-kit/permissions"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/state"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	state  *state.State
	db     *gorm.DB
	ginsta *ginsta.Ginsta
}

func (p *Plugin) Name() string {
	return "instagram"
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.state = params.State
	p.logger = params.Logger
	p.db = params.DB
	p.ginsta = ginsta.NewGinsta(
		&http.Client{
			Timeout: time.Second * 30,
		},
		nil,
	)

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

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name:        p.Name(),
		Description: "instagram.help.description",
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "instagram" &&
		event.Fields()[0] != "insta" {
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
