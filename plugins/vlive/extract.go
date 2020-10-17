package vlive

import (
	"regexp"
)

var vliveRegexp = regexp.MustCompile(".*vlive.tv/(channels/)?([^/]+).*")

func extractVLiveChannelID(input string) string {
	parts := vliveRegexp.FindStringSubmatch(input)
	if len(parts) < 3 {
		return input
	}

	return parts[2]
}
