package uploads

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (p *Plugin) Name() string {
	return "uploads"
}

func (p *Plugin) Start(params common.StartParameters) error {
	err := params.DB.AutoMigrate(Upload{}).Error
	if err != nil {
		return err
	}

	p.db = params.DB
	p.logger = params.Logger
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
		Description: "uploads.help.description",
		Commands: []common.Command{
			{
				Name:            "uploads.help.upload.name",
				Description:     "uploads.help.upload.description",
				SkipPrefix:      true,
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "@Bot", Type: common.Flag},
					{Name: "files…", Type: common.Attachment},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if event.Type == events.MessageCreateType && event.BotMention() && len(event.MessageCreate.Attachments) > 0 {
		p.handleUpload(event)
		return true
	}

	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "uploads" {
		return false
	}

	// TODO: Commands…

	p.handleList(event)
	return true
}
