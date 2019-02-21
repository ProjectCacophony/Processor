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

type sortByName []printValues

// Len is part of sort.Interface.
func (d sortByName) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d sortByName) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (d sortByName) Less(i, j int) bool {
	return strings.ToLower(d[i].Name) < strings.ToLower(d[j].Name)
}

func cmdElements(event *events.Event) {
	printMap := make(map[string][]printValues)

	for _, item := range list.TriggerList {
		printMap["triggers"] = append(printMap["triggers"], printValues{
			Name:        item.Name(),
			Description: event.Translate(item.Description()),
		})
	}
	for _, item := range list.FiltersList {
		printMap["filters"] = append(printMap["filters"], printValues{
			Name:        item.Name(),
			Description: event.Translate(item.Description()),
		})
	}
	for _, item := range list.ActionsList {
		printMap["actions"] = append(printMap["actions"], printValues{
			Name:        item.Name(),
			Description: event.Translate(item.Description()),
		})
	}

	sort.Sort(sortByName(printMap["triggers"]))
	sort.Sort(sortByName(printMap["filters"]))
	sort.Sort(sortByName(printMap["actions"]))

	_, err := event.Respond("automod.elements.response",
		"itemsMap", printMap,
	)
	event.Except(err)
}
