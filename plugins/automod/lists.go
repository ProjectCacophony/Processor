package automod

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/actions"
	"gitlab.com/Cacophony/Processor/plugins/automod/filters"
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/triggers"
)

// nolint: gochecknoglobals
var (
	triggerList = []interfaces.TriggerInterface{
		triggers.Message{},
	}

	filtersList = []interfaces.FilterInterface{
		filters.RegexMessageContent{},
	}

	actionsList = []interfaces.ActionInterface{
		actions.SendMessage{},
	}
)
