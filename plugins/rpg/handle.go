package rpg

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/Processor/plugins/rpg/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type Plugin struct{}

func (p *Plugin) Names() []string {
	return []string{"rpg"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	err := params.DB.AutoMigrate(models.Interaction{}).Error
	if err != nil {
		return err
	}

	repo := repo{db: params.DB}
	for _, happening := range happenings {
		happening.Init(&repo, params.State)
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
		Names: p.Names(),
		Hide:  true,
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		if event.Type == events.MessageCreateType &&
			!event.MessageCreate.Author.Bot {
			go p.processStep(event.GuildID, event.UserID, event.MessageCreate.Message)
			return false
		}

		return false
	}

	return false
}
