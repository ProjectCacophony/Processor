package serverlist

import (
	"regexp"
)

const (
	descriptionCharacterLimit = 100

	queueChannelIDKey = "cacophony:processor:serverlist:queue-channel-id"
	queueMessageKey   = "cacophony:processor:serverlist:queue-message"

	logChannelIDKey = "cacophony:processor:serverlist:log-channel-id"

	emojiApprove = "✅"
)

var (
	refreshQueueLock = func(guildID string) string {
		return "cacophony:processor:serverlist:queue-lock:guildid-" + guildID
	}

	redisListMessagesKey = func(channelID string) string {
		return "cacophony:processor:serverlist:list-messages:channelid-" + channelID
	}

	channelServersToPost = func(botID string) string {
		return "cacophony:processor:serverlist:channel-servers-to-post:botID-" + botID
	}

	serverNameInitialRegexp = regexp.MustCompile(`^[a-z]$`)
)
