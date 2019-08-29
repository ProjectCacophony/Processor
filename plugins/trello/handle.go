package trello

import (
	"errors"

	trello "github.com/VojtechVitek/go-trello"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	trello *trello.Client
	db     *gorm.DB
}

func (p *Plugin) Name() string {
	return "feedback"
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.db = params.DB

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil
	}

	if config.TrelloKey == "" || config.TrelloToken == "" {
		return errors.New("trello plugin configuration missing")
	}

	p.trello, err = trello.NewAuthClient(
		config.TrelloKey,
		&config.TrelloToken,
	)
	if err != nil {
		return errors.New("trello plugin unable to initialize client")
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

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name:        p.Name(),
		Description: "trello.help.description",
		Commands: []common.Command{
			{
				Name:            "trello.help.suggest.name",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "suggest", Type: common.Flag},
					{Name: "suggestion title", Type: common.QuotedText},
					{Name: "suggestion description", Type: common.QuotedText},
				},
			},
			{
				Name:            "trello.help.issue.name",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "issue", Type: common.Flag},
					{Name: "issue title", Type: common.QuotedText},
					{Name: "issue description", Type: common.QuotedText},
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
	case "suggest", "suggestion", "issue", "bug":
		event.Require(func() {
			p.handleSuggestion(event)
		}, permissions.BotAdmin)
		return true
	}

	return false
}
