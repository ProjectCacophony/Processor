package gall

import (
	"net/http"
	"strings"
	"time"

	"gitlab.com/Cacophony/go-kit/permissions"

	"github.com/Seklfreak/ginside"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/state"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
	"go.uber.org/zap"
)

type Plugin struct {
	logger  *zap.Logger
	state   *state.State
	db      *gorm.DB
	ginside *ginside.GInside
}

func (p *Plugin) Name() string {
	return "gall"
}

func (p *Plugin) Start(params common.StartParameters) error {
	var err error

	p.state = params.State
	p.logger = params.Logger
	p.db = params.DB
	p.ginside = ginside.NewGInside(&http.Client{
		Timeout: time.Minute,
	})

	err = params.DB.AutoMigrate(Entry{}).Error
	if err != nil {
		return err
	}

	return nil
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

func (p *Plugin) Localisations() []interfaces.Localisation {
	local, err := localisation.NewFileSource("assets/translations/gall.en.toml", "en")
	if err != nil {
		p.logger.Error("failed to load localisation", zap.Error(err))
	}

	return []interfaces.Localisation{local}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "gall" {
		return false
	}

	if len(event.Fields()) > 1 {
		if strings.EqualFold(event.Fields()[1], "remove") {
			event.Require(func() {
				p.remove(event)
			}, permissions.DiscordManageChannels)

			return true
		}

		event.Require(func() {
			p.add(event)
		}, permissions.DiscordManageChannels,
		)
		return true
	}

	p.status(event)

	return true
}
