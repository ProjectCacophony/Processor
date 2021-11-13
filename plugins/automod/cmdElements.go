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
	triggers := make([]printValues, 0, len(list.TriggerList))
	for _, item := range list.TriggerList {
		if item.Deprecated() {
			continue
		}

		triggers = append(triggers, printValues{
			Name:        item.Name(),
			Description: event.Translate(item.Description()),
		})
	}
	filters := make([]printValues, 0, len(list.FiltersList))
	for _, item := range list.FiltersList {
		if item.Deprecated() {
			continue
		}

		filters = append(filters, printValues{
			Name:        item.Name(),
			Description: event.Translate(item.Description()),
		})
	}
	actions := make([]printValues, 0, len(list.ActionsList))
	for _, item := range list.ActionsList {
		if item.Deprecated() {
			continue
		}

		actions = append(actions, printValues{
			Name:        item.Name(),
			Description: event.Translate(item.Description()),
		})
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
