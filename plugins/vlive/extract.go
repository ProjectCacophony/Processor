package vlive

import (
	"regexp"
)

var vliveRegexp = regexp.MustCompile(".*vlive.(tv|com)/(channels/)?([^/]+).*")

func extractVLiveChannelID(input string) string {
	parts := vliveRegexp.FindStringSubmatch(input)
	if len(parts) < 4 {
		return input
	}

	return parts[3]
}
