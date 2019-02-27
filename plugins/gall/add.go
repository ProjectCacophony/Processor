package gall

import (
	"strings"

	"github.com/Seklfreak/ginside"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) add(event *events.Event) {
	fields := event.Fields()[1:]

	channel, fields, err := paramsExtractChannel(event, fields)
	if err != nil {
		event.Except(err)
		return
	}

	all, fields := paramsIsAll(fields)

	if len(fields) < 1 {
		event.Respond("gall.add.too-few") // nolint: errcheck
		return
	}

	boardID := strings.ToLower(fields[0])

	entries, err := entryFindMany(p.db,
		"guild_id = ?",
		event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}
	if len(entries) >= feedsPerGuildLimit(event) &&
		feedsPerGuildLimit(event) >= 0 {
		event.Respond("gall.add.too-many") // nolint: errcheck
		return
	}

	for _, entry := range entries {
		if entry.ChannelID != channel.ID {
			continue
		}
		if entry.BoardID != boardID {
			continue
		}

		event.Respond("gall.add.duplicate") // nolint: errcheck
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
				event.Respond("gall.add.not-found") // nolint: errcheck
				return
			}
			event.Except(err)
			return
		}
	}

	if len(posts) == 0 {
		event.Respond("gall.add.not-found") // nolint: errcheck
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
	)
	if err != nil {
		event.Except(err)
	}

	_, err = event.Respond("gall.add.success",
		"boardID", boardID,
		"channel", channel,
		"recommended", !all,
	)
	event.Except(err)
}
