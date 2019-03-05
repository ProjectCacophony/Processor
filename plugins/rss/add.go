package rss

import (
	"net/url"
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) add(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("rss.add.too-few") // nolint: errcheck
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
		event.Respond("rss.add.too-few") // nolint: errcheck
		return
	}

	inputURL := strings.Trim(fields[0], "<>")
	feedURL := inputURL

	_, err = url.ParseRequestURI(feedURL)
	if err != nil {
		event.Respond("rss.add.not-found") // nolint: errcheck
		return
	}

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
		event.Respond("rss.add.too-many") // nolint: errcheck
		return
	}

	feed, err := getFeed(p.httpClient, p.parser, feedURL)
	if err != nil || len(feed.Items) == 0 {
		feedURL, err = getFeedURLFromPage(p.httpClient, feedURL)
		if err != nil {
			event.Respond("rss.add.not-found") // nolint: errcheck
			return
		}

		feed, err = getFeed(p.httpClient, p.parser, feedURL)
		if err != nil || len(feed.Items) == 0 {
			event.Respond("rss.add.not-found") // nolint: errcheck
			return
		}
	}

	for _, entry := range entries {
		if entry.ChannelID != channel.ID {
			continue
		}
		if !strings.EqualFold(entry.FeedURL, feed.FeedLink) {
			continue
		}

		event.Respond("rss.add.duplicate") // nolint: errcheck
		return
	}

	if feed.FeedLink == "" {
		feed.FeedLink = feedURL
	}
	if feed.Link == "" {
		feed.Link = inputURL
	}

	err = entryAdd(
		p.db,
		event.UserID,
		channel.ID,
		event.GuildID,
		feed.Title,
		feed.Link,
		feed.FeedLink,
		event.BotUserID,
		event.DM(),
	)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("rss.add.success",
		"feed", feed,
		"channel", channel,
		"dm", event.DM(),
	)
	event.Except(err)
}
