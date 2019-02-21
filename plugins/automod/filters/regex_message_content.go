package filters

import (
	"regexp"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/go-kit/events"
)

type RegexMessageContent struct {
}

func (f RegexMessageContent) Name() string {
	return "if_message_content_regex"
}

func (f RegexMessageContent) NewItem(env *models.Env, value string) (interfaces.FilterItemInterface, error) {
	regex, err := regexp.Compile(value)
	if err != nil {
		return nil, err
	}

	return &RegexMessageContentItem{
		regexp: regex,
	}, nil
}

type RegexMessageContentItem struct {
	regexp *regexp.Regexp
}

func (f *RegexMessageContentItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.MessageCreateType {
		return false
	}

	return f.regexp.MatchString(env.Event.MessageCreate.Content)
}
