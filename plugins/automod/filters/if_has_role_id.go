package filters

import (
	"errors"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
)

// deprecated: use if_role
type RoleID struct {
}

func (f RoleID) Name() string {
	return "if_has_role_id"
}

func (f RoleID) Args() int {
	return 1
}

func (f RoleID) Deprecated() bool {
	return true
}

func (f RoleID) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	RoleIDs := strings.Split(args[0], ",")

	if len(RoleIDs) == 0 {
		return nil, errors.New("no Role IDs defined")
	}

	return &RoleIDItem{
		RoleIDs: RoleIDs,
	}, nil
}

func (f RoleID) Description() string {
	return "automod.filters.if_has_role_id"
}

type RoleIDItem struct {
	RoleIDs []string
}

func (f *RoleIDItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	for _, userID := range env.UserID {
		member, err := env.State.Member(env.GuildID, userID)
		if err != nil {
			return false
		}

		for _, roleID := range f.RoleIDs {
			if matchChannels(member.Roles, roleID) {
				continue
			}

			return false
		}
	}

	return true
}
