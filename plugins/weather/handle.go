package weather

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"github.com/shawntoffel/darksky"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type Plugin struct {
	logger  *zap.Logger
	db      *gorm.DB
	darkSky darksky.DarkSky
	config  Config
	redis   *redis.Client
}

func (p *Plugin) Name() string {
	return "weather"
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.db = params.DB
	p.redis = params.Redis

	err := envconfig.Process("", &p.config)
	if err != nil {
		return nil
	}

	if p.config.GoogleMapsKey == "" || p.config.DarkSkyKey == "" {
		return errors.New("weather plugin configuration missing")
	}

	p.darkSky = darksky.New(p.config.DarkSkyKey)

	return params.DB.AutoMigrate(Weather{}).Error
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
		Description: "weather.help.description",
		Commands: []common.Command{
			{
				Name:            "weather.help.view",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "weather", Type: common.Flag},
				},
			},
			{
				Name:            "weather.help.set",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "weather", Type: common.Flag},
					{Name: "set", Type: common.Flag},
					{Name: "location", Type: common.Text},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	switch event.Fields()[0] {
	case "weather":

		if len(event.Fields()) > 2 {
			switch event.Fields()[1] {
			case "set":
				p.setUserWeather(event)
				return true
			}
		}

		p.viewWeather(event)
		return true
	}

	return false
}
