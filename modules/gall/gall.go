package gall

import (
	"time"

	"strings"

	"context"

	"github.com/Seklfreak/ginside"
	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo/bson"
	"github.com/opentracing/opentracing-go"
	"gitlab.com/Cacophony/SqsProcessor/models"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/mdb"
	"gitlab.com/Cacophony/dhelpers/state"
)

func displayBoard(ctx context.Context, event dhelpers.EventContainer) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "gall.displayBoard")
	defer span.Finish()

	// get boardID
	boardID := strings.ToLower(event.Args[1])

	// start typing
	event.GoType(event.MessageCreate.ChannelID)

	// get all posts (if requested), or just recommended (default)
	recommended := true
	if includes, _ := dhelpers.SliceContainsLowerExclude(event.Args, []string{"all"}); includes {
		recommended = false
	}

	// get data
	var err error
	var posts []ginside.Post
	posts, err = ginside.BoardPosts(boardID, recommended)
	dhelpers.CheckErr(err)

	// if no posts found, try if it's a minor board
	if len(posts) <= 0 {
		posts, err = ginside.BoardMinorPosts(boardID, recommended)
		dhelpers.CheckErr(err)
	}

	// if still no posts found, post error and stop
	if len(posts) <= 0 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, "GallNotFound")
		dhelpers.CheckErr(err)
		return
	}

	// build embed
	embed := &discordgo.MessageEmbed{
		URL:   friendlyBoardURL(boardID),
		Title: dhelpers.Tf("GallBoardPostsTitle", "boardID", boardID, "recommended", recommended),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    dhelpers.T("GallEmbedFooter"),
			IconURL: GallIcon,
		},
		Color: GallColor,
	}

	// build embed description
	for i, post := range posts {
		embed.Description += dhelpers.Tf("GallPostSummary", "post", post) + "\n"
		// only show top 10
		if i >= 9 {
			break
		}
	}

	// send away
	_, err = event.SendEmbed(event.MessageCreate.ChannelID, embed)
	dhelpers.CheckErr(err)
}

func addBoard(ctx context.Context, event dhelpers.EventContainer) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "gall.addBoard")
	defer span.Finish()

	// we need at least three args
	if len(event.Args) < 3 {
		return
	}

	// start typing
	event.GoType(event.MessageCreate.ChannelID)

	// get all posts, or just recommended (default)
	recommended := true
	var includes bool
	if includes, event.Args = dhelpers.SliceContainsLowerExclude(event.Args, []string{"all"}); includes {
		recommended = false
	}

	// get boardID from args
	boardID := strings.ToLower(event.Args[2])
	var minorGallery bool

	// get data
	posts, err := ginside.BoardPosts(boardID, recommended)
	dhelpers.CheckErr(err)

	// if no posts found, try minor board
	if len(posts) <= 0 {
		minorGallery = true
		posts, err = ginside.BoardMinorPosts(boardID, recommended)
		dhelpers.CheckErr(err)
	}

	// if still no posts, post error and stop
	if len(posts) <= 0 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, "GallNotFound")
		dhelpers.CheckErr(err)
		return
	}

	// target channel = source channel, or different channel if specified
	var targetChannel *discordgo.Channel
	targetChannel, err = state.Channel(event.MessageCreate.ChannelID)
	dhelpers.CheckErr(err)
	if len(event.Args) >= 4 {
		targetChannel, err = state.ChannelFromMention(event.MessageCreate.GuildID, event.Args[3])
		dhelpers.CheckErr(err)
	}

	// check if we already have this board set up
	if alreadySetUp(boardID, targetChannel.ID) {
		_, err = event.SendMessagef(event.MessageCreate.ChannelID, "GallBoardFeedAlreadySetUp", "boardID", boardID, "targetChannel", targetChannel)
		dhelpers.CheckErr(err)
		return
	}

	// create gall feed entry
	entry := models.GallFeedEntry{
		GuildID:       targetChannel.GuildID,
		ChannelID:     targetChannel.ID,
		AddedByUserID: event.MessageCreate.Author.ID,
		BoardID:       boardID,
		MinorGallery:  minorGallery,
		Recommended:   recommended,
		LastCheck:     time.Now(),
		AddedAt:       time.Now(),
	}

	// insert gall feed entry into db
	_, err = mdb.Insert(models.GallTable, entry)
	dhelpers.CheckErr(err)

	// send success message
	_, err = event.SendMessagef(event.MessageCreate.ChannelID, "GallBoardFeedAdded", "entry", entry, "targetChannel", targetChannel)
	dhelpers.CheckErr(err)
}

func listFeeds(ctx context.Context, event dhelpers.EventContainer) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "gall.listFeeds")
	defer span.Finish()

	// get all gall feed entries on current server
	var err error
	var feedEntries []models.GallFeedEntry
	err = mdb.Iter(models.GallTable.DB().Find(bson.M{"guildid": event.MessageCreate.GuildID})).All(&feedEntries)
	dhelpers.CheckErr(err)

	// if no entries found, post error and stop
	if len(feedEntries) <= 0 {
		_, err = event.SendMessagef(event.MessageCreate.ChannelID, "GallNoFeeds")
		dhelpers.CheckErr(err)
		return
	}

	// create gall feed entries list message
	var message string
	for _, feedEntry := range feedEntries {
		message += dhelpers.Tf("GallFeedEntry", "feedEntry", feedEntry) + "\n"
	}
	message += event.Tf("GallFeedEntriesSummary", "feedEntryCount", len(feedEntries))

	// send away
	_, err = event.SendMessage(event.MessageCreate.ChannelID, message)
	dhelpers.CheckErr(err)
}

func removeFeed(ctx context.Context, event dhelpers.EventContainer) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "gall.removeFeed")
	defer span.Finish()

	var err error

	// we need at least three args
	if len(event.Args) < 3 {
		return
	}

	// get board ID from args
	boardID := strings.ToLower(event.Args[2])

	// try finding gall feed entries with the given boardID on the current server
	var feedEntries []models.GallFeedEntry
	err = mdb.Iter(models.GallTable.DB().Find(bson.M{
		"boardid": boardID,
		"guildid": event.MessageCreate.GuildID,
	})).All(&feedEntries)
	// if none found, post error and stop
	if len(feedEntries) <= 0 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, "GallFeedNotFound")
		dhelpers.CheckErr(err)
		return
	}
	dhelpers.CheckErr(err)

	// figure out which one to delete
	toDelete := feedEntries[0]
	// delete in current channel first
	for _, entry := range feedEntries {
		if entry.ChannelID == event.MessageCreate.ChannelID {
			toDelete = entry
			break
		}
	}

	// delete entry from DB
	err = mdb.DeleteID(models.GallTable, toDelete.ID)
	dhelpers.CheckErr(err)

	// send success message
	_, err = event.SendMessagef(event.MessageCreate.ChannelID, "GallFeedEntryRemoved", "feedEntry", toDelete)
	dhelpers.CheckErr(err)
}
