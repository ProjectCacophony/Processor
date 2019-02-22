package filters

import (
	"errors"
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

func (f RegexMessageContent) Args() int {
	return 1
}

func (f RegexMessageContent) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	regex, err := regexp.Compile(args[0])
	if err != nil {
		return nil, err
	}

	return &RegexMessageContentItem{
		regexp: regex,
	}, nil
}

func (f RegexMessageContent) Description() string {
	return "automod.filters.if_message_content_regex"
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
