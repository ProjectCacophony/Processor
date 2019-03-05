// nolint: dupl
package filters

import (
	"errors"
	"regexp"
	"strconv"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

// nolint: gochecknoglobals
var emojiRegex = regexp.MustCompile(`[\x{00A0}-\x{1F9EF}]|<(a)?:[^<>:]+:[0-9]+>`)

type EmojiCount struct {
}

func (f EmojiCount) Name() string {
	return "if_emoji_count"
}

func (f EmojiCount) Args() int {
	return 2
}

func (f EmojiCount) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
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

	return &EmojiCountItem{
		Amount:     amount,
		Comparison: comparisonType,
	}, nil
}

func (f EmojiCount) Description() string {
	return "automod.filters.if_emoji_count"
}

type EmojiCountItem struct {
	Amount     int
	Comparison AmountComparisonType
}

func (f *EmojiCountItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.MessageCreateType {
		return false
	}

	amount := len(emojiRegex.FindAllString(env.Event.MessageCreate.Content, -1))

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
