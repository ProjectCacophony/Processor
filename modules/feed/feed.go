package feed

import (
	"strings"

	"net/url"

	"time"

	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/Cacophony/SqsProcessor/models"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/mdb"
	"gitlab.com/Cacophony/dhelpers/state"
)

func displayFeed(event dhelpers.EventContainer) {
	event.GoType(event.MessageCreate.ChannelID)

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
		if dhelpers.IsNetworkErr(err) {
			_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNetworkErr", "feedURL", feedURL)
			dhelpers.CheckErr(err)
			return
		}
		if strings.Contains(err.Error(), "Failed to detect feed type") {
			// if no feed read, try to read feed url from html page
			var feedURLNew string
			feedURLNew, err = getFeedURLFromPage(feedURL)
			if err != nil {
				if strings.Contains(err.Error(), "unable to find feed url") {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNotFound", "feedURL", feedURL)
					dhelpers.CheckErr(err)
					return
				}
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
				if strings.Contains(err.Error(), "Failed to detect feed type") {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNotFound", "feedURL", feedURLNew)
					dhelpers.CheckErr(err)
					return
				}
				if dhelpers.IsNetworkErr(err) {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNetworkErr", "feedURL", feedURLNew)
					dhelpers.CheckErr(err)
					return
				}
			}

		}
	}
	dhelpers.CheckErr(err)

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

	if feed.Items != nil {
		for i, item := range feed.Items {
			embed.Description += "\n" + dhelpers.Tf("FeedItemSummary", "item", item)
			// display last five
			if i >= 4 {
				break
			}
		}
	}

	if feed.PublishedParsed != nil && !feed.PublishedParsed.IsZero() {
		embed.Footer.Text += "| Updated at "
		embed.Timestamp = dhelpers.DiscordTime(*feed.PublishedParsed)
	}

	if feed.Image != nil {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: feed.Image.URL,
		}
	}

	_, err = event.SendEmbed(event.MessageCreate.ChannelID, embed)
	dhelpers.CheckErr(err)
}

func addFeed(event dhelpers.EventContainer) {
	if len(event.Args) < 3 {
		return
	}

	event.GoType(event.MessageCreate.ChannelID)

	feedURL := dhelpers.CleanURL(event.Args[2])

	sourceChannel, err := state.Channel(event.MessageCreate.ChannelID)
	dhelpers.CheckErr(err)

	// check url is valid
	if _, err := url.ParseRequestURI(feedURL); err != nil {
		_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedInvalidURL", "feedURL", feedURL)
		dhelpers.CheckErr(err)
		return
	}

	// try to read feed
	feed, err := GetFeed(feedURL)
	if err != nil {
		if dhelpers.IsNetworkErr(err) {
			_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNetworkErr", "feedURL", feedURL)
			dhelpers.CheckErr(err)
			return
		}
		if strings.Contains(err.Error(), "Failed to detect feed type") {
			// if no feed read, try to read feed url from html page
			var feedURLNew string
			feedURLNew, err = getFeedURLFromPage(feedURL)
			if err != nil {
				if strings.Contains(err.Error(), "unable to find feed url") {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNotFound", "feedURL", feedURL)
					dhelpers.CheckErr(err)
					return
				}
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
				if strings.Contains(err.Error(), "Failed to detect feed type") {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNotFound", "feedURL", feedURLNew)
					dhelpers.CheckErr(err)
					return
				}
				if dhelpers.IsNetworkErr(err) {
					_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNetworkErr", "feedURL", feedURLNew)
					dhelpers.CheckErr(err)
					return
				}
			}

		}
	}
	dhelpers.CheckErr(err)

	targetChannel := sourceChannel
	if len(event.Args) >= 4 {
		targetChannel, err = state.ChannelFromMention(sourceChannel.GuildID, event.Args[3])
		dhelpers.CheckErr(err)
	}

	if alreadySetUp(feedURL, targetChannel.ID) {
		_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedAlreadySetUp", "feedURL", feedURL, "targetChannel", targetChannel)
		dhelpers.CheckErr(err)
		return
	}

	_, err = mdb.Insert(models.FeedTable, models.FeedEntry{
		GuildID:       targetChannel.GuildID,
		ChannelID:     targetChannel.ID,
		AddedByUserID: event.MessageCreate.Author.ID,
		LastCheck:     time.Now(),
		FeedURL:       feedURL,
		FeedTitle:     feed.Title,
	})
	dhelpers.CheckErr(err)

	_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedAdded", "feed", feed, "feedURL", feedURL, "targetChannel", targetChannel)
	dhelpers.CheckErr(err)
}

func listFeeds(event dhelpers.EventContainer) {
	sourceChannel, err := state.Channel(event.MessageCreate.ChannelID)
	dhelpers.CheckErr(err)

	var feedEntries []models.FeedEntry
	err = mdb.Iter(models.FeedTable.DB().Find(bson.M{"guildid": sourceChannel.GuildID})).All(&feedEntries)
	dhelpers.CheckErr(err)

	if len(feedEntries) <= 0 {
		_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedNoFeeds")
		dhelpers.CheckErr(err)
		return
	}

	var message string
	for _, entry := range feedEntries {
		message += dhelpers.Tf("FeedEntry", "entry", entry) + "\n"
	}
	message += event.Tf("FeedEntriesSummary", "feedEntryCount", len(feedEntries))

	_, err = event.SendMessage(event.MessageCreate.ChannelID, message)
	dhelpers.CheckErr(err)
}

func removeFeed(event dhelpers.EventContainer) {
	sourceChannel, err := state.Channel(event.MessageCreate.ChannelID)
	dhelpers.CheckErr(err)

	if len(event.Args) < 3 {
		return
	}

	feedURL := strings.ToLower(event.Args[2])

	var feedEntries []models.FeedEntry
	err = mdb.Iter(models.FeedTable.DB().Find(bson.M{
		"feedurl": bson.M{"$regex": bson.RegEx{Pattern: "^" + regexp.QuoteMeta(feedURL) + "$", Options: "i"}},
		"guildid": sourceChannel.GuildID,
	})).All(&feedEntries)
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
		if entry.ChannelID == sourceChannel.ID {
			toDelete = entry
			break
		}
	}

	err = mdb.DeleteID(models.FeedTable, toDelete.ID)
	dhelpers.CheckErr(err)

	_, err = event.SendMessagef(event.MessageCreate.ChannelID, "FeedEntryRemoved", "entry", toDelete)
	dhelpers.CheckErr(err)
}
