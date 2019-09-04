package uploads

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
	"mvdan.cc/xurls/v2"
)

var xurlsStrict = xurls.Strict()

type Plugin struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (p *Plugin) Names() []string {
	return []string{"uploads"}
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
		Names:       p.Names(),
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
			{
				Name:        "uploads.help.list.name",
				Description: "uploads.help.list.description",
			},
			{
				Name:        "uploads.help.delete.name",
				Description: "uploads.help.delete.description",
				Params: []common.CommandParam{
					{Name: "delete", Type: common.Flag},
					{Name: "link", Type: common.Link},
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

	if xurlsStrict.MatchString(event.OriginalCommand()) {
		p.handleUpload(event)
		return true
	}

	if event.Fields()[0] != "uploads" {
		return false
	}

	if len(event.Fields()) >= 2 && (event.Fields()[1] == "delete" || event.Fields()[1] == "remove") {
		p.handleDelete(event)
		return true
	}

	p.handleList(event)
	return true
}
