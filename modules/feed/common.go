package feed

import (
	"bytes"

	"strings"

	"net/url"

	"regexp"

	"strconv"
	"time"

	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/Processor/models"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/net"
)

// GetFeed returns the gofeed.Feed for an URL (ATOM or RSS)
func GetFeed(feedURL string) (feed *gofeed.Feed, err error) {
	parsedFeedURL, err := url.Parse(feedURL)
	if err != nil {
		return nil, err
	}

	// add cache busting
	newQueries := parsedFeedURL.Query()
	newQueries.Set("_", strconv.FormatInt(time.Now().Unix(), 10))
	parsedFeedURL.RawQuery = newQueries.Encode()

	// download feed page
	var feedData []byte
	feedData, err = net.Get(parsedFeedURL.String())
	if err != nil {
		return nil, err
	}

	// create a new feed parser and try parsing the feed
	feedParser := gofeed.NewParser()
	feed, err = feedParser.Parse(bytes.NewReader(feedData))
	return feed, err
}

// getFeedURLFromPage tries to find a feedURL on a webpage (using HTML headers)
func getFeedURLFromPage(pageURL string) (feedURL string, err error) {
	// check the pageURL by trying to parse it
	var parsedURL *url.URL
	parsedURL, err = url.Parse(pageURL)
	if err != nil {
		return "", err
	}

	// download web page
	var pageData []byte
	pageData, err = net.Get(pageURL)
	if err != nil {
		return "", err
	}

	// create a new goquery document from the web page content
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(pageData))
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

// alreadySetUp returns true if the feed is already set up in the channel
func alreadySetUp(ctx context.Context, feedURL, channelID string) (already bool) {
	// go to database, check if feedURL (casesensitive) is already set up for the given channel
	count, err := models.FeedRepository.Count(ctx, bson.NewDocument(
		bson.EC.Regex("feedurl", "^"+regexp.QuoteMeta(feedURL)+"$", "i"),
		bson.EC.String("channelid", channelID),
	))
	dhelpers.LogError(err)
	// return true if one or more results
	return count > 0
}

// GetEntryID returns the ID for a Feed Entry (used for deduplication)
func GetEntryID(entry *gofeed.Item) string {
	// try GUID if not empty
	if entry.GUID != "" {
		return entry.GUID
	}
	// fallback to link
	return entry.Link
}
