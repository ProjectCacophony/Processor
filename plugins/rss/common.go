package rss

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/PuerkitoBio/goquery"

	"gitlab.com/Cacophony/go-kit/permissions"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func feedsPerGuildLimit(event *events.Event) int {
	if event.Has(permissions.BotOwner) {
		return -1
	}

	if event.DM() {
		return 2
	}

	return 5
}

func paramsExtractChannel(event *events.Event, args []string) (*discordgo.Channel, []string, error) {
	for i, arg := range args {
		channel, err := event.State().ChannelFromMention(event.GuildID, arg)
		if err != nil {
			continue
		}

		return channel, append(args[:i], args[i+1:]...), nil
	}

	channel, err := event.State().Channel(event.ChannelID)
	return channel, args, err
}

func getFeed(client *http.Client, parser *gofeed.Parser, feedURL string) (*gofeed.Feed, error) {
	parsedFeedURL, err := url.Parse(feedURL)
	if err != nil {
		return nil, err
	}

	// add cache busting
	newQueries := parsedFeedURL.Query()
	newQueries.Set("_", strconv.FormatInt(time.Now().Unix(), 10))
	parsedFeedURL.RawQuery = newQueries.Encode()

	// download feed page
	resp, err := client.Get(parsedFeedURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parser.Parse(resp.Body)
}

func getFeedURLFromPage(client *http.Client, pageURL string) (string, error) {
	// check the pageURL by trying to parse it
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		return "", err
	}

	// download web page
	resp, err := client.Get(pageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// create a new goquery document from the web page content
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// go through all links
	links := doc.Find("link")
	for _, link := range links.Nodes {
		// check link time, if feed link type return that one
		entryNode := goquery.NewDocumentFromNode(link)
		linkType, _ := entryNode.Attr("type")
		switch linkType {
		case "application/atom+xml":
			linkHref, _ := entryNode.Attr("href")
			if linkHref == "" {
				continue
			}
			// make url absolute
			if !strings.HasPrefix(linkHref, "http") {
				linkHref = parsedURL.Scheme + "://" + parsedURL.Host + linkHref
			}
			return linkHref, nil
		case "application/rss+xml":
			linkHref, _ := entryNode.Attr("href")
			if linkHref == "" {
				continue
			}
			// make url absolute
			if !strings.HasPrefix(linkHref, "http") {
				linkHref = parsedURL.Scheme + "://" + parsedURL.Host + linkHref
			}
			return linkHref, nil
		}
	}

	// we failed to find feed url, return error
	return "", errors.New("unable to find feed url")
}
