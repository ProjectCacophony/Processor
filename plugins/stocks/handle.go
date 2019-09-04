package stocks

import (
	"errors"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/external/iexcloud"
	"go.uber.org/zap"
)

type Plugin struct {
	logger    *zap.Logger
	iexClient *iexcloud.IEX
	db        *gorm.DB
}

func (p *Plugin) Names() []string {
	return []string{"stocks", "stock"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.db = params.DB

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil
	}

	if config.IEXAPISecret == "" {
		return errors.New("stocks plugin configuration missing")
	}

	p.iexClient = iexcloud.NewIEX(
		&http.Client{
			Timeout: time.Second * 10,
		},
		config.IEXAPISecret,
	)

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

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Names:       p.Names(),
		Description: "stocks.help.description",
		Commands: []common.Command{
			{
				Name:        "stocks.help.display.name",
				Description: "stocks.help.display.description",
				Params: []common.CommandParam{
					{Name: "Symbol", Type: common.Text},
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
	case "stocks", "stock":

		p.handleStocks(event)
		return true
	}

	return false
}
