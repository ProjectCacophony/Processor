package serverlist

const (
	descriptionCharacterLimit = 100

	queueChannelIDKey = "cacophony:processor:serverlist:queue-channel-id"
	queueMessageKey   = "cacophony:processor:serverlist:queue-message"
)

// nolint: gochecknoglobals
var (
	refreshQueueLock = func(guildID string) string {
		return "cacophony:processor:serverlist:queue-lock:guildid-" + guildID
	}
)
