package gall

import (
	"context"

	"github.com/Seklfreak/ginside"
	"github.com/mongodb/mongo-go-driver/bson"
	"gitlab.com/Cacophony/Processor/models"
	"gitlab.com/Cacophony/dhelpers"
)

var (
	// GallColor is the discord colour used for Gall embeds
	GallColor = dhelpers.HexToDecimal("#4064D0")
	// GallIcon is a Gall Logo
	GallIcon = "https://i.imgur.com/tIYs6jt.png"

	friendlyBoardURL = func(boardID string) string {
		return "http://gall.dcinside.com/board/lists/?id=" + boardID + "&page=1&exception_mode=recommend"
	}
)

// alreadySetUp returns true if the board is already set up in the channel
func alreadySetUp(ctx context.Context, boardID, channelID string) (already bool) {
	count, err := models.GallRepository.Count(
		ctx,
		bson.NewDocument(
			bson.EC.String("boardid", boardID),
			bson.EC.String("channelid", channelID),
		),
	)
	dhelpers.LogError(err)
	return count > 0
}

// GetEntryID returns the ID for a Feed Entry (used for deduplication)
func GetEntryID(entry ginside.Post) string {
	return entry.URL
}
