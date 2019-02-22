package automod

import (
	"sort"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/list"
	"gitlab.com/Cacophony/go-kit/events"
)

type printValues struct {
	Name        string
	Description string
}

type sortPrintValuesByName []printValues

// Len is part of sort.Interface
func (d sortPrintValuesByName) Len() int {
	return len(d)
}

// Swap is part of sort.Interface
func (d sortPrintValuesByName) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface
func (d sortPrintValuesByName) Less(i, j int) bool {
	return strings.ToLower(d[i].Name) < strings.ToLower(d[j].Name)
}

func (p *Plugin) cmdElements(event *events.Event) {

	triggers := make([]printValues, len(list.TriggerList))
	for i, item := range list.TriggerList {
		triggers[i] = printValues{
			Name:        item.Name(),
			Description: event.Translate(item.Description()),
		}
	}
	filters := make([]printValues, len(list.FiltersList))
	for i, item := range list.FiltersList {
		filters[i] = printValues{
			Name:        item.Name(),
			Description: event.Translate(item.Description()),
		}
	}
	actions := make([]printValues, len(list.ActionsList))
	for i, item := range list.ActionsList {
		actions[i] = printValues{
			Name:        item.Name(),
			Description: event.Translate(item.Description()),
		}
	}

	sort.Sort(sortPrintValuesByName(triggers))
	sort.Sort(sortPrintValuesByName(filters))
	sort.Sort(sortPrintValuesByName(actions))

	_, err := event.Respond("automod.elements.response",
		"triggers", triggers,
		"filters", filters,
		"actions", actions,
	)
	event.Except(err)
}
