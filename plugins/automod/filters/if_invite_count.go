package filters

import (
	"errors"
	"regexp"
	"strconv"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

var discordInviteRegex = regexp.MustCompile(
	`(http(s)?:\/\/)?(discord\.gg(\/invite)?|discordapp\.com\/invite)\/([A-Za-z0-9-]+)`,
)

type InvitesCount struct {
}

func (f InvitesCount) Name() string {
	return "if_invites_count"
}

func (f InvitesCount) Args() int {
	return 2
}

func (f InvitesCount) Deprecated() bool {
	return false
}

func (f InvitesCount) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	comparisonType, err := extractAmountComparisonType(args[0])
	if err != nil {
		return nil, err
	}

	amount, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}

	if amount < 0 {
		return nil, errors.New("amount has to be 0 or greater")
	}

	return &InvitesCountItem{
		Amount:     amount,
		Comparison: comparisonType,
	}, nil
}

func (f InvitesCount) Description() string {
	return "automod.filters.if_invites_count"
}

type InvitesCountItem struct {
	Amount     int
	Comparison AmountComparisonType
}

func (f *InvitesCountItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.MessageCreateType {
		return false
	}

	amount := len(discordInviteRegex.FindAllString(env.Event.MessageCreate.Content, -1))

	switch f.Comparison {
	case AmountComparisonLT:
		return amount < f.Amount
	case AmountComparisonEQ:
		return amount == f.Amount
	case AmountComparisonGT:
		return amount > f.Amount
	}

	return false
}
