package weverse

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) add(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("weverse.add.too-few", "communities", p.communities)
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
		event.Respond("weverse.add.too-few", "communities", p.communities)
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
		event.Respond("weverse.add.too-many")
		return
	}

	community, err := extractCommunity(p.communities, fields)
	if err != nil {
		if strings.Contains(err.Error(), "community not found") {
			event.Respond("weverse.add.not-found", "communities", p.communities)
			return
		}
		event.Except(err)
		return
	}

	for _, entry := range entries {
		if entry.ChannelOrUserID != channel.ID {
			continue
		}
		if entry.WeverseChannelID != community.ID {
			continue
		}

		event.Respond("weverse.add.duplicate")
		return
	}

	err = entryAdd(
		p.db,
		event.UserID,
		channel.ID,
		event.GuildID,
		community.Name,
		community.ID,
		event.BotUserID,
		event.DM(),
	)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("weverse.add.success",
		"community", community,
		"channel", channel,
		"dm", event.DM(),
	)
	event.Except(err)
}
