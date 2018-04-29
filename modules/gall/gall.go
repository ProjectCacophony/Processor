package gall

import (
	"time"

	"regexp"

	"github.com/Seklfreak/ginside"
	"github.com/bwmarrin/discordgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/project-d-collab/SqsProcessor/models"
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/mdb"
	"gitlab.com/project-d-collab/dhelpers/state"
)

func displayBoard(event dhelpers.EventContainer) {
	boardID := event.Args[1]

	event.GoType(event.MessageCreate.ChannelID)

	// get data
	posts, err := ginside.BoardRecommendedPosts(boardID)
	dhelpers.CheckErr(err)

	if len(posts) <= 0 {
		posts, err = ginside.BoardMinorRecommendedPosts(boardID)
		dhelpers.CheckErr(err)
	}

	if len(posts) <= 0 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, "GallNotFound")
		dhelpers.CheckErr(err)
		return
	}

	// build embed
	embed := &discordgo.MessageEmbed{
		URL:   friendlyBoardURL(boardID),
		Title: dhelpers.Tf("GallBoardPostsTitle", "boardID", boardID),
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

func addBoard(event dhelpers.EventContainer) {
	if len(event.Args) < 3 {
		return
	}

	event.GoType(event.MessageCreate.ChannelID)

	sourceChannel, err := state.Channel(event.MessageCreate.ChannelID)
	dhelpers.CheckErr(err)

	boardID := event.Args[2]
	var minorGallery bool

	// get data
	posts, err := ginside.BoardRecommendedPosts(boardID)
	dhelpers.CheckErr(err)

	if len(posts) <= 0 {
		minorGallery = true
		posts, err = ginside.BoardMinorRecommendedPosts(boardID)
		dhelpers.CheckErr(err)
	}

	if len(posts) <= 0 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, "GallNotFound")
		dhelpers.CheckErr(err)
		return
	}

	targetChannel := sourceChannel
	if len(event.Args) >= 4 {
		targetChannel, err = state.ChannelFromMention(sourceChannel.GuildID, event.Args[3])
		dhelpers.CheckErr(err)
	}

	_, err = mdb.Insert(models.GallTable, models.GallFeedEntry{
		GuildID:       targetChannel.GuildID,
		ChannelID:     targetChannel.ID,
		AddedByUserID: event.MessageCreate.Author.ID,
		BoardID:       boardID,
		LastCheck:     time.Now(),
		MinorGallery:  minorGallery,
	})
	dhelpers.CheckErr(err)

	_, err = event.SendMessagef(event.MessageCreate.ChannelID, "GallBoardFeedAdded", "boardID", boardID, "targetChannel", targetChannel)
	dhelpers.CheckErr(err)
}

func listFeeds(event dhelpers.EventContainer) {
	sourceChannel, err := state.Channel(event.MessageCreate.ChannelID)
	dhelpers.CheckErr(err)

	var feedEntries []models.GallFeedEntry
	err = mdb.Iter(models.GallTable.DB().Find(bson.M{"guildid": sourceChannel.GuildID})).All(&feedEntries)
	dhelpers.CheckErr(err)

	if len(feedEntries) <= 0 {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, "GallNoFeeds")
		dhelpers.CheckErr(err)
		return
	}

	var message string
	for _, feedEntry := range feedEntries {
		message += dhelpers.Tf("GallFeedEntry", "feedEntry", feedEntry) + "\n"
	}
	message += dhelpers.Tf("GallFeedEntriesSummary", "feedEntryCount", len(feedEntries))

	_, err = event.SendMessage(event.MessageCreate.ChannelID, message)
	dhelpers.CheckErr(err)
}

func removeFeed(event dhelpers.EventContainer) {
	sourceChannel, err := state.Channel(event.MessageCreate.ChannelID)
	dhelpers.CheckErr(err)

	if len(event.Args) < 3 {
		return
	}

	boardID := event.Args[2]

	var feedEntry models.GallFeedEntry
	err = mdb.One(
		models.GallTable.DB().Find(
			bson.M{
				"boardid": bson.M{"$regex": bson.RegEx{Pattern: "^" + regexp.QuoteMeta(boardID) + "$", Options: "i"}},
				"guildid": sourceChannel.GuildID,
			}),
		&feedEntry,
	)
	if mdb.ErrNotFound(err) {
		_, err = event.SendMessage(event.MessageCreate.ChannelID, "GallFeedNotFound")
		dhelpers.CheckErr(err)
		return
	}
	dhelpers.CheckErr(err)

	err = mdb.DeleteID(models.GallTable, feedEntry.ID)
	dhelpers.CheckErr(err)

	_, err = event.SendMessagef(event.MessageCreate.ChannelID, "GallFeedEntryRemoved", "feedEntry", feedEntry)
	dhelpers.CheckErr(err)
}
