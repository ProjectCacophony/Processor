package vlive

import (
	"regexp"
)

var (
	vliveRegexp          = regexp.MustCompile(".*vlive.(tv|com)/(channels/)?([^/]+).*")
	vliveChannelIDRegexp = regexp.MustCompile("^[A-Z0-9]+$")
)

func extractVLiveChannelID(input string) string {
	parts := vliveRegexp.FindStringSubmatch(input)
	if len(parts) < 4 {
		return input
	}

	if !vliveChannelIDRegexp.MatchString(parts[3]) {
		return ""
	}

	return parts[3]
}
