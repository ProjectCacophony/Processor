package gall

import (
	"strings"

	"github.com/Seklfreak/ginside"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) add(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("gall.add.too-few")
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

	all, fields := paramsIsAll(fields)

	if len(fields) < 1 {
		event.Respond("gall.add.too-few")
		return
	}

	boardID := strings.ToLower(fields[0])

	entries, err := entryFindMany(p.db,
		"((guild_id = ? AND dm = false) OR (channel_id = ? AND dm = true)) AND dm = ?",
		event.GuildID, event.UserID, event.DM(),
	)
	if err != nil {
		event.Except(err)
		return
	}
	if len(entries) >= feedsPerGuildLimit(event) &&
		feedsPerGuildLimit(event) >= 0 {
		event.Respond("gall.add.too-many")
		return
	}

	for _, entry := range entries {
		if entry.ChannelID != channel.ID {
			continue
		}
		if entry.BoardID != boardID {
			continue
		}

		event.Respond("gall.add.duplicate")
		return
	}

	var posts []ginside.Post
	var minorGallery bool

	posts, err = p.ginside.BoardPosts(event.Context(), boardID, false)
	if err != nil || len(posts) == 0 {
		minorGallery = true
		posts, err = p.ginside.BoardMinorPosts(event.Context(), boardID, false)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				event.Respond("gall.add.not-found")
				return
			}
			event.Except(err)
			return
		}
	}

	if len(posts) == 0 {
		event.Respond("gall.add.not-found")
		return
	}

	err = entryAdd(
		p.db,
		event.UserID,
		channel.ID,
		event.GuildID,
		boardID,
		minorGallery,
		!all,
		event.BotUserID,
		event.DM(),
	)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("gall.add.success",
		"boardID", boardID,
		"channel", channel,
		"recommended", !all,
		"dm", event.DM(),
	)
	event.Except(err)
}
