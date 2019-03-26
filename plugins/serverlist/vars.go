package serverlist

import (
	"regexp"
)

const (
	descriptionCharacterLimit = 100

	queueChannelIDKey = "cacophony:processor:serverlist:queue-channel-id"
	queueMessageKey   = "cacophony:processor:serverlist:queue-message"

	emojiApprove = "✅"
	emojiReject  = "❌"
)

// nolint: gochecknoglobals
var (
	refreshQueueLock = func(guildID string) string {
		return "cacophony:processor:serverlist:queue-lock:guildid-" + guildID
	}

	serverNameInitialRegexp = regexp.MustCompile(`^[a-z]$`)
)
