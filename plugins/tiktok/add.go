package tiktok

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) add(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("tiktok.add.too-few")
		return
	}

	fields := event.Fields()[2:]

	channel, fields, err := paramsExtractChannel(event, fields)
	if err != nil {
		event.Except(err)
		return
	}
	if event.DM() {
		channel.ID = event.UserID
	}

	if len(fields) < 1 {
		event.Respond("tiktok.add.too-few")
		return
	}

	entries, err := entryFindMany(p.db,
		"((guild_id = ? AND dm = false) OR (channel_or_user_id = ? AND dm = true)) AND dm = ?",
		event.GuildID, event.UserID, event.DM(),
	)
	if err != nil {
		event.Except(err)
		return
	}
	if len(entries) >= feedsPerGuildLimit(event) &&
		feedsPerGuildLimit(event) >= 0 {
		event.Respond("tiktok.add.too-many")
		return
	}

	username := extractTikTokUsername(fields[0])
	if len(username) <= 0 {
		event.Respond("tiktok.add.not-found")
		return
	}

	for _, entry := range entries {
		if entry.ChannelOrUserID != channel.ID {
			continue
		}
		if !strings.EqualFold(entry.TikTokUsername, username) {
			continue
		}

		event.Respond("tiktok.add.duplicate")
		return
	}

	err = entryAdd(
		p.db,
		event.UserID,
		channel.ID,
		event.GuildID,
		username,
		event.BotUserID,
		event.DM(),
	)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("tiktok.add.success",
		"tiktokUsername", username,
		"channel", channel,
		"dm", event.DM(),
	)
	event.Except(err)
}
