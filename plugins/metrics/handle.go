package metrics

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
)

type Plugin struct {
}

func (p *Plugin) Names() []string {
	return []string{"metrics"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	err := params.DB.AutoMigrate(Counter{}).Error
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		err = metric.Register(params.DB)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) Stop(params common.StopParameters) error {
	return nil
}

func (p *Plugin) Priority() int {
	return 10000
}

func (p *Plugin) Passthrough() bool {
	return true
}

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Names: p.Names(),
		Hide:  true,
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if event.Command() {
		err := totalCommands.Inc(event.DB())
		event.ExceptSilent(err)
	}

	if event.Type == events.MessageCreateType {
		err := totalMessagesReceived.Inc(event.DB())
		event.ExceptSilent(err)
	}

	return false
}
