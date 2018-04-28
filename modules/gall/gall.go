package gall

import (
	"strings"

	"github.com/Seklfreak/ginside"
	"github.com/bwmarrin/discordgo"
	"gitlab.com/project-d-collab/dhelpers"
)

func displayBoard(event dhelpers.EventContainer) {
	boardID := event.Args[1]

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
		// gall https page is broken
		post.Link = strings.Replace(post.Link, "https://", "http://", 1)
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
