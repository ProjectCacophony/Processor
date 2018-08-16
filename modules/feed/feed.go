package feed

import (
	"strings"

	"net/url"

	"time"

	"regexp"

	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo/bson"
	"github.com/mmcdole/gofeed"
	"github.com/opentracing/opentracing-go"
	"gitlab.com/Cacophony/SqsProcessor/models"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/mdb"
	"gitlab.com/Cacophony/dhelpers/state"
)

func displayFeed(ctx context.Context, event dhelpers.EventContainer) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "feed.displayFeed")
	defer span.Finish()

	// start typing
	event.GoType()

	// clean url from discord
	feedURL := dhelpers.CleanURL(event.Args[1])

	// check url is valid
	if _, err := url.ParseRequestURI(feedURL); err != nil {
		_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedInvalidURL", "feedURL", feedURL)
		dhelpers.CheckErr(err)
		return
	}

	// try to read feed
	feed, err := GetFeed(feedURL)
	if err != nil {
		// if network error post error and stop
		if dhelpers.IsNetworkErr(err) {
			_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNetworkErr", "feedURL", feedURL)
			dhelpers.CheckErr(err)
			return
		}
		// if no feed read, try to read feed url from html page
		if strings.Contains(err.Error(), "Failed to detect feed type") {
			// try find feed url
			var feedURLNew string
			feedURLNew, err = getFeedURLFromPage(feedURL)
			if err != nil {
				// if no feed url found, post error and stop
				if strings.Contains(err.Error(), "unable to find feed url") {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNotFound", "feedURL", feedURL)
					dhelpers.CheckErr(err)
					return
				}
				// if network error post error and stop
				if dhelpers.IsNetworkErr(err) {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNetworkErr", "feedURL", feedURL)
					dhelpers.CheckErr(err)
					return
				}
			}
			dhelpers.CheckErr(err)
			// try to read new feed url
			feed, err = GetFeed(feedURLNew)
			if err != nil {
				// if no feed read post error and stop
				if strings.Contains(err.Error(), "Failed to detect feed type") {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNotFound", "feedURL", feedURLNew)
					dhelpers.CheckErr(err)
					return
				}
				// if network error post error and stop
				if dhelpers.IsNetworkErr(err) {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNetworkErr", "feedURL", feedURLNew)
					dhelpers.CheckErr(err)
					return
				}
			}

		}
	}
	dhelpers.CheckErr(err)

	// set feed link to human friendly link if possible
	feedLink := feedURL
	if feed.Link != "" {
		feedLink = feed.Link
	}

	// build embed
	embed := &discordgo.MessageEmbed{
		URL:         feedLink,
		Title:       feed.Title,
		Description: feed.Description + "\n",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "URL: " + feedURL,
		},
	}

	// display items if possible
	if feed.Items != nil {
		for i, item := range feed.Items {
			embed.Description += "\n" + dhelpers.Tf("FeedItemSummary", "item", item)
			// display last five
			if i >= 4 {
				break
			}
		}
	}

	// display last feed update if possible
	if feed.PublishedParsed != nil && !feed.PublishedParsed.IsZero() {
		embed.Footer.Text += "| Updated at "
		embed.Timestamp = dhelpers.DiscordTime(*feed.PublishedParsed)
	}

	// add thumbnail if poossible
	if feed.Image != nil {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: feed.Image.URL,
		}
	}

	// send
	_, err = event.SendEmbed(event.MessageCreate.ChannelID, embed)
	dhelpers.CheckErr(err)
}

func addFeed(ctx context.Context, event dhelpers.EventContainer) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "feed.addFeed")
	defer span.Finish()

	// we need at least three args
	if len(event.Args) < 3 {
		return
	}

	// start typing
	event.GoType()

	feedURL := dhelpers.CleanURL(event.Args[2])

	sourceChannel, err := state.Channel(event.MessageCreate.ChannelID)
	dhelpers.CheckErr(err)

	// check url is valid
	_, err = url.ParseRequestURI(feedURL)
	if err != nil {
		_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedInvalidURL", "feedURL", feedURL)
		dhelpers.CheckErr(err)
		return
	}

	// try to read feed
	var feed *gofeed.Feed
	feed, err = GetFeed(feedURL)
	if err != nil {
		// if network error post error and stop
		if dhelpers.IsNetworkErr(err) {
			_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNetworkErr", "feedURL", feedURL)
			dhelpers.CheckErr(err)
			return
		}
		// if no feed read, try to read feed url from html page
		if strings.Contains(err.Error(), "Failed to detect feed type") {
			// try find feed url
			var feedURLNew string
			feedURLNew, err = getFeedURLFromPage(feedURL)
			if err != nil {
				// if no feed url found, post error and stop
				if strings.Contains(err.Error(), "unable to find feed url") {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNotFound", "feedURL", feedURL)
					dhelpers.CheckErr(err)
					return
				}
				// if network error post error and stop
				if dhelpers.IsNetworkErr(err) {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNetworkErr", "feedURL", feedURL)
					dhelpers.CheckErr(err)
					return
				}
			}
			dhelpers.CheckErr(err)
			// try to read new feed url
			feed, err = GetFeed(feedURLNew)
			if err != nil {
				// if no feed read post error and stop
				if strings.Contains(err.Error(), "Failed to detect feed type") {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNotFound", "feedURL", feedURLNew)
					dhelpers.CheckErr(err)
					return
				}
				// if network error post error and stop
				if dhelpers.IsNetworkErr(err) {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNetworkErr", "feedURL", feedURLNew)
					dhelpers.CheckErr(err)
					return
				}
			}

		}
	}
	dhelpers.CheckErr(err)

	// target channel = source channel, overwrite with channel from arguments if given
	targetChannel := sourceChannel
	if len(event.Args) >= 4 {
		targetChannel, err = state.ChannelFromMention(event.MessageCreate.GuildID, event.Args[3])
		dhelpers.CheckErr(err)
	}

	// post error if we have the feed set up in the target channel ready
	if alreadySetUp(feedURL, targetChannel.ID) {
		_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedAlreadySetUp", "feedURL", feedURL, "targetChannel", targetChannel)
		dhelpers.CheckErr(err)
		return
	}

	// insert new entry
	_, err = mdb.Insert(models.FeedTable, models.FeedEntry{
		GuildID:       targetChannel.GuildID,
		ChannelID:     targetChannel.ID,
		AddedByUserID: event.MessageCreate.Author.ID,
		LastCheck:     time.Now(),
		FeedURL:       feedURL,
		FeedTitle:     feed.Title,
		AddedAt:       time.Now(),
	})
	dhelpers.CheckErr(err)

	// send success message
	_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedAdded", "feed", feed, "feedURL", feedURL, "targetChannel", targetChannel)
	dhelpers.CheckErr(err)
}

func listFeeds(ctx context.Context, event dhelpers.EventContainer) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "feed.listFeeds")
	defer span.Finish()

	// request feeds on this guild
	var err error
	var feedEntries []models.FeedEntry
	err = mdb.Iter(models.FeedTable.DB().Find(bson.M{"guildid": event.MessageCreate.GuildID})).All(&feedEntries)
	dhelpers.CheckErr(err)

	// post error if no feed set up so far
	if len(feedEntries) <= 0 {
		_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNoFeeds")
		dhelpers.CheckErr(err)
		return
	}

	// create feed list message
	var message string
	for _, entry := range feedEntries {
		message += dhelpers.Tf("FeedEntry", "entry", entry) + "\n"
	}
	message += event.Tf("FeedEntriesSummary", "feedEntryCount", len(feedEntries))

	// send to discord
	_, err = event.SendMessage(event.MessageCreate.ChannelID, message)
	dhelpers.CheckErr(err)
}

func removeFeed(ctx context.Context, event dhelpers.EventContainer) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "feed.removeFeed")
	defer span.Finish()

	var err error

	// we need at least three args
	if len(event.Args) < 3 {
		return
	}

	// get feed URL to delete from args
	feedURL := strings.ToLower(event.Args[2])

	// try to find feed with given feedURL (case insensitive) on current server
	var feedEntries []models.FeedEntry
	err = mdb.Iter(models.FeedTable.DB().Find(bson.M{
		"feedurl": bson.M{"$regex": bson.RegEx{Pattern: "^" + regexp.QuoteMeta(feedURL) + "$", Options: "i"}},
		"guildid": event.MessageCreate.GuildID,
	})).All(&feedEntries)
	// if not found, post error and stop
	if len(feedEntries) <= 0 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, "FeedEntryNotFound")
		dhelpers.CheckErr(err)
		return
	}
	dhelpers.CheckErr(err)

	// figure out which one to delete
	toDelete := feedEntries[0]
	// delete in current channel first
	for _, entry := range feedEntries {
		if entry.ChannelID == event.MessageCreate.GuildID {
			toDelete = entry
			break
		}
	}

	// delete entry
	err = mdb.DeleteID(models.FeedTable, toDelete.ID)
	dhelpers.CheckErr(err)

	// send success message
	_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedEntryRemoved", "entry", toDelete)
	dhelpers.CheckErr(err)
}
