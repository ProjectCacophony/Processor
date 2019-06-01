package instagram

import (
	"regexp"
)

var instagramRegex = regexp.MustCompile(".*instagram.com/([^/]+).*")

func extractInstagramUsername(input string) string {
	parts := instagramRegex.FindStringSubmatch(input)
	if len(parts) < 2 {
		return input
	}

	return parts[1]
}
