package filters

import (
	"errors"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
)

type User struct {
}

func (f User) Name() string {
	return "if_is_user"
}

func (f User) Args() int {
	return 1
}

func (f User) Deprecated() bool {
	return false
}

func (f User) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	var userIDs []string
	for _, mention := range strings.Split(args[0], ",") {
		user, err := env.State.UserFromMention(mention)
		if err != nil {
			continue
		}

		userIDs = append(userIDs, user.ID)
	}

	if len(userIDs) == 0 {
		return nil, errors.New("no users found")
	}

	return &UserItem{
		UserIDs: userIDs,
	}, nil
}

func (f User) Description() string {
	return "automod.filters.if_is_user"
}

type UserItem struct {
	UserIDs []string
}

func (f *UserItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	for _, userID := range env.UserID {
		if sliceContains(f.UserIDs, userID) {
			continue
		}

		return false
	}

	return true
}
