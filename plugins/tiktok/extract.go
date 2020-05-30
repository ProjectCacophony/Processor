package tiktok

import (
	"regexp"
)

var (
	usernameLinkRegex = regexp.MustCompile(".*tiktok.com/@([^/]+).*")
	usernameRegex     = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
)

func extractTikTokUsername(input string) string {
	candidate := input

	parts := usernameLinkRegex.FindStringSubmatch(input)
	if len(parts) >= 2 {
		candidate = parts[1]
	}

	if usernameRegex.MatchString(candidate) {
		return candidate
	}
	return ""
}
