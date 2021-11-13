package filters

import (
	"errors"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

// deprecated: use if_is_user
type UserID struct{}

func (f UserID) Name() string {
	return "if_is_user_id"
}

func (f UserID) Args() int {
	return 1
}

func (f UserID) Deprecated() bool {
	return true
}

func (f UserID) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	UserIDs := strings.Split(args[0], ",")

	if len(UserIDs) == 0 {
		return nil, errors.New("no User IDs defined")
	}

	return &UserIDItem{
		UserIDs: UserIDs,
	}, nil
}

func (f UserID) Description() string {
	return "automod.filters.if_is_user_id"
}

type UserIDItem struct {
	UserIDs []string
}

func (f *UserIDItem) Match(env *models.Env) bool {
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
