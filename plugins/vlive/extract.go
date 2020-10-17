package vlive

import (
	"regexp"
)

var vliveRegexp = regexp.MustCompile(".*vlive.tv/([^/]+).*")

func extractVLiveChannelID(input string) string {
	parts := vliveRegexp.FindStringSubmatch(input)
	if len(parts) < 2 {
		return input
	}

	return parts[1]
}
