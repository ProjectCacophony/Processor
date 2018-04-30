package feed

import (
	"strings"

	"net/url"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/project-d-collab/dhelpers"
)

func displayFeed(event dhelpers.EventContainer) {
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
			feedURL, err = getFeedURLFromPage(feedURL)
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
			feed, err = GetFeed(feedURL)
			if err != nil {
				if strings.Contains(err.Error(), "Failed to detect feed type") {
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
