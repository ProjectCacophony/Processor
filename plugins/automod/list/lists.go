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
	}

	FiltersList = []interfaces.FilterInterface{
		filters.RegexMessageContent{},
		filters.True{},
	}

	ActionsList = []interfaces.ActionInterface{
		actions.SendMessage{},
		actions.ApplyRole{},
		actions.IncrBucket{},
	}
)
