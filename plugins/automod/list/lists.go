package list

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/actions"
	"gitlab.com/Cacophony/Processor/plugins/automod/filters"
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/triggers"
)

// nolint: gochecknoglobals
var (
	TriggerList = []interfaces.TriggerInterface{
		triggers.Message{},
		triggers.BucketUpdated{},
		triggers.Join{},
	}

	FiltersList = []interfaces.FilterInterface{
		filters.RegexMessageContent{},
		filters.True{},
		filters.BucketGT{},
		filters.RegexUserName{},
		filters.AccountAgeLT{},
		filters.MentionsCount{},
	}

	ActionsList = []interfaces.ActionInterface{
		actions.SendMessage{},
		actions.ApplyRole{},
		actions.IncrBucket{},
		actions.SendMessageTo{},
	}
)
