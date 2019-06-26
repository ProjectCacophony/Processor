package quickactions

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type Plugin struct {
	logger    *zap.Logger
	redis     *redis.Client
	publisher *events.Publisher
	tokens    map[string]string
	db        *gorm.DB
}

func (p *Plugin) Name() string {
	return "quickactions"
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.redis = params.Redis
	p.publisher = params.Publisher
	p.tokens = params.Tokens
	p.db = params.DB

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
		Reactions: []common.Reaction{
			{
				EmojiName:   "<:quickaction_star:579322267100512256>", // TODO: do not hardcode emoji ids
				Description: "quickactions.help.star.description",
			},
			{
				EmojiName:   "<:quickaction_remind_1h:579342313495592980>",
				Description: "quickactions.help.remind1.description",
			},
			{
				EmojiName:   "<:quickaction_remind_8h:579342223372582912>",
				Description: "quickactions.help.remind8.description",
			},
			{
				EmojiName:   "<:quickaction_remind_24h:579342223141896252>",
				Description: "quickactions.help.remind24.description",
			},
			{
				EmojiName:   "<:quickaction_remind_custom:582444736519077889>",
				Description: "quickactions.help.remindCustom.description",
			},
			{
				EmojiName:   "<:quickaction_raw:593498575275556875>",
				Description: "quickactions.help.raw.description",
			},
		},
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

	if event.Type == events.CacophonyQuestionnaireMatch &&
		event.QuestionnaireMatch.Key == questionnaireRemindKey {
		p.handleRemindQuestionnaire(event)
		return true
	}

	return false
}
