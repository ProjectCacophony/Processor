package gall

import (
	"time"

	"github.com/Seklfreak/ginside"
	"github.com/bwmarrin/discordgo"
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
			IconURL: dhelpers.T("GallEmbedIcon"),
		},
		Color: color,
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

	// get data
	posts, err := ginside.BoardRecommendedPosts(boardID)
	dhelpers.CheckErr(err)

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
	})
	dhelpers.CheckErr(err)

	_, err = event.SendMessagef(event.MessageCreate.ChannelID, "GallBoardFeedAdded", "boardID", boardID, "targetChannel", targetChannel)
	dhelpers.CheckErr(err)
}
