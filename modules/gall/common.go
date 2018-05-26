package gall

import (
	"github.com/Seklfreak/ginside"
	"github.com/globalsign/mgo/bson"
	"gitlab.com/Cacophony/SqsProcessor/models"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/mdb"
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
func alreadySetUp(boardID, channelID string) (already bool) {
	count, _ := mdb.Count(
		models.GallTable, bson.M{
			"boardid":   boardID,
			"channelid": channelID,
		})
	return count > 0
}

// GetEntryID returns the ID for a Feed Entry (used for deduplication)
func GetEntryID(entry ginside.Post) string {
	if entry.ID != "" {
		return entry.ID
	}
	return entry.URL
}
