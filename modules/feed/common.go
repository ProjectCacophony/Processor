package feed

import (
	"bytes"

	"strings"

	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"gitlab.com/project-d-collab/dhelpers"
)

// GetFeed returns the gofeed.Feed for an URL (ATOM or RSS)
func GetFeed(feedURL string) (feed *gofeed.Feed, err error) {
	var feedData []byte
	feedData, err = dhelpers.NetGet(feedURL)
	if err != nil {
		return nil, err
	}

	feedParser := gofeed.NewParser()

	feed, err = feedParser.Parse(bytes.NewReader(feedData))
	return feed, err
}

// getFeedURLFromPage tries to find a feedURL on a webpage (using HTML headers)
func getFeedURLFromPage(pageURL string) (feedURL string, err error) {
	var pageData []byte
	pageData, err = dhelpers.NetGet(pageURL)
	if err != nil {
		return "", err
	}

	var parsedURL *url.URL
	parsedURL, err = url.Parse(pageURL)
	if err != nil {
		return "", err
	}

	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(pageData))
	if err != nil {
		return "", err
	}

	links := doc.Find("link")
	for _, link := range links.Nodes {
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

	return "", errors.New("unable to find feed url")
}
