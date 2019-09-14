package tools

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
	"mvdan.cc/xurls/v2"
)

var xurlsStrict = xurls.Strict()

func (p *Plugin) handleShorten(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("tools.shorten.too-few")
		return
	}

	var link string
	var candidate string

	for _, field := range event.Fields() {
		candidate = strings.Trim(field, "<>")
		if xurlsStrict.MatchString(candidate) {
			link = candidate
		}
	}

	if link == "" {
		event.Respond("tools.shorten.too-few")
		return
	}

	_, err := event.Respond("tools.shorten.content", "link", link)
	event.Except(err)
}
