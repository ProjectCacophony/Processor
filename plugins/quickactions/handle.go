package quickactions

import (
	"github.com/go-redis/redis"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type Plugin struct {
	logger    *zap.Logger
	redis     *redis.Client
	publisher *events.Publisher
	tokens    map[string]string
}

func (p *Plugin) Name() string {
	return "quickactions"
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.redis = params.Redis
	p.publisher = params.Publisher
	p.tokens = params.Tokens

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
		Description: "quickactions.help.description",
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if event.Type == events.MessageReactionAddType {
		return p.handleReaction(event)
	}

	if event.Type == events.CacophonyQuickactionRemind {
		p.handleQuickactionRemind(event)
		return true
	}

	if !event.Command() {
		return false
	}

	return false
}
