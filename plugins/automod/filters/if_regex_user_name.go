package filters

import (
	"errors"
	"regexp"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type RegexUserName struct{}

func (f RegexUserName) Name() string {
	return "if_user_name_regex"
}

func (f RegexUserName) Args() int {
	return 1
}

func (f RegexUserName) Deprecated() bool {
	return false
}

func (f RegexUserName) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	regex, err := regexp.Compile(args[0])
	if err != nil {
		return nil, err
	}

	return &RegexUserNameItem{
		regexp: regex,
	}, nil
}

func (f RegexUserName) Description() string {
	return "automod.filters.if_user_name_regex"
}

type RegexUserNameItem struct {
	regexp *regexp.Regexp
}

func (f *RegexUserNameItem) Match(env *models.Env) bool {
	var didNotMatch bool

	for _, userID := range env.UserID {
		user, err := env.State.User(userID)
		if err != nil {
			return false
		}

		if f.regexp.MatchString(user.Username) {
			continue
		}

		didNotMatch = true
	}

	return !didNotMatch
}
